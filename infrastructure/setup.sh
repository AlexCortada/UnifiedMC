#!/bin/bash
# =====================================================
# Unified IT Operations Portal - Ubuntu Setup Script
# Run this on a fresh Ubuntu 24.04 LTS installation
# Usage: chmod +x setup.sh && sudo ./setup.sh
# =====================================================

set -e

echo "=========================================="
echo "  Unified IT Operations Portal - Setup"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# =====================================================
# Step 1: System Prerequisites
# =====================================================
echo ""
log_info "Step 1: Installing system prerequisites..."

export DEBIAN_FRONTEND=noninteractive

# Update system
apt-get update -y
apt-get upgrade -y

# Install essential tools
apt-get install -y \
    curl \
    wget \
    git \
    vim \
    htop \
    net-tools \
    ca-certificates \
    gnupg \
    lsb-release \
    jq \
    unzip \
    postgresql-client \
    ufw

log_info "System prerequisites installed."

# =====================================================
# Step 2: Install Docker
# =====================================================
echo ""
log_info "Step 2: Installing Docker..."

# Add Docker's official GPG key
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

# Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Start and enable Docker
systemctl enable docker
systemctl start docker

# Add current user to docker group (requires re-login to take effect)
usermod -aG docker $USER 2>/dev/null || true

log_info "Docker installed: $(docker --version)"

# =====================================================
# Step 3: Install Go
# =====================================================
echo ""
log_info "Step 3: Installing Go 1.22..."

GO_VERSION="1.22.5"
wget -q "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -O /tmp/go.tar.gz
tar -C /usr/local -xzf /tmp/go.tar.gz
rm /tmp/go.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bashrc
export PATH=$PATH:/usr/local/go/bin

# Create GOPATH
mkdir -p /root/go/{bin,src,pkg}
echo 'export GOPATH=/root/go' >> /root/.bashrc
echo 'export PATH=$PATH:/root/go/bin' >> /root/.bashrc

log_info "Go installed: $(go version)"

# =====================================================
# Step 4: Clone the Repository
# =====================================================
echo ""
log_info "Step 4: Cloning UnifiedMC repository..."

APP_DIR="/opt/unifiedmc"

if [ -d "$APP_DIR" ]; then
    log_warn "Directory $APP_DIR already exists. Pulling latest..."
    cd $APP_DIR
    git pull origin main
else
    git clone https://github.com/AlexCortada/UnifiedMC.git $APP_DIR
    cd $APP_DIR
fi

log_info "Repository ready at $APP_DIR"

# =====================================================
# Step 5: Configure System for Services
# =====================================================
echo ""
log_info "Step 5: Configuring system for services..."

# Elasticsearch requires vm.max_map_count >= 262144
if [ -f /etc/sysctl.conf ]; then
    if ! grep -q "vm.max_map_count" /etc/sysctl.conf; then
        echo "vm.max_map_count=262144" >> /etc/sysctl.conf
        sysctl -p 2>/dev/null || true
    fi
fi

# Increase file descriptor limits
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

log_info "System configuration complete."

# =====================================================
# Step 6: Start Infrastructure Services
# =====================================================
echo ""
log_info "Step 6: Starting infrastructure services..."

cd $APP_DIR

# Remove obsolete version field if present
if head -1 infrastructure/docker/docker-compose.yml | grep -q "^version:"; then
    sed -i '/^version:/d' infrastructure/docker/docker-compose.yml
    # Remove empty line at top if present
    sed -i '/^$/N;/^\n$/d' infrastructure/docker/docker-compose.yml
fi

# Start all services
docker compose -f infrastructure/docker/docker-compose.yml up -d

log_info "All services starting. This may take a minute..."

# Wait for services to be healthy
echo ""
log_info "Waiting for services to become healthy..."

MAX_WAIT=120
ELAPSED=0

while [ $ELAPSED -lt $MAX_WAIT ]; do
    # Check if postgres is responding
    if PGPASSWORD=*** psql -h localhost -U unifiedmc -d unifiedmc -c "SELECT 1" &>/dev/null; then
        log_info "PostgreSQL is healthy!"
        break
    fi
    sleep 5
    ELAPSED=$((ELAPSED + 5))
    echo -n "."
done

if [ $ELAPSED -ge $MAX_WAIT ]; then
    log_warn "Timed out waiting for PostgreSQL. Services may still be starting..."
fi

# Give extra time for other services
sleep 10

# =====================================================
# Step 7: Run Database Migrations
# =====================================================
echo ""
log_info "Step 7: Running database migrations..."

for f in $APP_DIR/database/migrations/*.sql; do
    log_info "  Running: $(basename $f)"
    PGPASSWORD=*** psql -h localhost -U unifiedmc -d unifiedmc -f "$f" 2>/dev/null || true
done

log_info "Database migrations complete."

# =====================================================
# Step 8: Build and Start API Gateway
# =====================================================
echo ""
log_info "Step 8: Building API Gateway..."

cd $APP_DIR/services/api-gateway
export PATH=$PATH:/usr/local/go/bin
go mod tidy
go build -o api-gateway .

log_info "API Gateway built successfully."

# Start API gateway in background
nohup ./api-gateway > /var/log/uop-api.log 2>&1 &
echo $! > /var/log/uop-api.pid

sleep 2

# Test the API
if curl -s http://localhost:8080/health | grep -q "healthy"; then
    log_info "API Gateway is running and healthy!"
else
    log_warn "API Gateway may still be starting. Check with: curl http://localhost:8080/health"
fi

# =====================================================
# Step 9: Configure Firewall
# =====================================================
echo ""
log_info "Step 9: Configuring firewall..."

ufw --force reset
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp      # SSH
ufw allow 8080/tcp    # API Gateway
ufw allow 5432/tcp    # PostgreSQL
ufw allow 6379/tcp    # Redis
ufw allow 9200/tcp    # Elasticsearch
ufw allow 8123/tcp    # ClickHouse HTTP
ufw allow 9000/tcp    # ClickHouse Native
ufw allow 9092/tcp    # Kafka
ufw --force enable

log_info "Firewall configured."

# =====================================================
# Step 10: Create Systemd Service for API Gateway
# =====================================================
echo ""
log_info "Step 10: Creating systemd service for API Gateway..."

cat > /etc/systemd/system/uop-api.service << 'SYSTEMD'
[Unit]
Description=Unified IT Operations Portal - API Gateway
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/unifiedmc/services/api-gateway
ExecStart=/opt/unifiedmc/services/api-gateway/api-gateway
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=SERVICE_NAME=uop-api-gateway

[Install]
WantedBy=multi-user.target
SYSTEMD

systemctl daemon-reload
systemctl enable uop-api
log_info "Systemd service created. Use: systemctl start uop-api"

# =====================================================
# Step 11: Create Convenience Scripts
# =====================================================
echo ""
log_info "Step 11: Creating convenience scripts..."

mkdir -p /opt/unifiedmc/scripts

# Start all services
cat > /opt/unifiedmc/scripts/start.sh << 'SCRIPT'
#!/bin/bash
cd /opt/unifiedmc
docker compose -f infrastructure/docker/docker-compose.yml up -d
echo "Services starting. Check with: docker ps"
SCRIPT

# Stop all services
cat > /opt/unifiedmc/scripts/stop.sh << 'SCRIPT'
#!/bin/bash
cd /opt/unifiedmc
docker compose -f infrastructure/docker/docker-compose.yml down
echo "Services stopped."
SCRIPT

# Status check
cat > /opt/unifiedmc/scripts/status.sh << 'SCRIPT'
#!/bin/bash
echo "=== Container Status ==="
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
echo ""
echo "=== API Health ==="
curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || echo "API not responding"
echo ""
echo "=== Database ==="
PGPASSWORD=*** psql -h localhost -U unifiedmc -d unifiedmc -c "SELECT count(*) as device_count FROM unified_devices;" 2>/dev/null || echo "Database not ready"
SCRIPT

# View logs
cat > /opt/unifiedmc/scripts/logs.sh << 'SCRIPT'
#!/bin/bash
SERVICE=${1:-""}
if [ -z "$SERVICE" ]; then
    echo "Usage: ./logs.sh <service-name>"
    echo "Available: postgres, redis, kafka, zookeeper, elasticsearch, clickhouse"
    exit 1
fi
docker logs -f docker-${SERVICE}-1
SCRIPT

chmod +x /opt/unifiedmc/scripts/*.sh

log_info "Scripts created in /opt/unifiedmc/scripts/"

# =====================================================
# Final Summary
# =====================================================
echo ""
echo "=========================================="
echo "  Setup Complete!"
echo "=========================================="
echo ""
echo "  Application Directory: $APP_DIR"
echo "  API Gateway:           http://localhost:8080"
echo "  Health Check:          http://localhost:8080/health"
echo "  Device List:           http://localhost:8080/api/v1/devices"
echo ""
echo "  Services Running:"
docker ps --format "    {{.Names}}: {{.Status}}" 2>/dev/null || echo "    (checking...)"
echo ""
echo "  Convenience Scripts:"
echo "    ./scripts/start.sh    - Start all services"
echo "    ./scripts/stop.sh     - Stop all services"
echo "    ./scripts/status.sh   - Check status"
echo "    ./logs.sh <service>   - View logs"
echo ""
echo "  Next Steps:"
echo "    1. Test the API: curl http://localhost:8080/health"
echo "    2. View devices: curl http://localhost:8080/api/v1/devices"
echo "    3. Start building connectors (Intune, Ivanti, Entra)"
echo ""
echo "=========================================="

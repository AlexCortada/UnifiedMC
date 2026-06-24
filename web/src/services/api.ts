export interface Device {
  id: string;
  external_id: string;
  connector_type: string;
  tenant_id: string;
  canonical_name: string;
  asset_type: string;
  os_type: string;
  os_version: string;
  serial_number: string;
  compliance_status: string;
  status: string;
  last_seen: string;
}

export interface DashboardSummary {
  total_devices: { value: number; previous_value: number; change_percent: number };
  online_devices: { value: number; percentage: number };
  offline_devices: { value: number; percentage: number };
  compliance_rate: { rate: number; compliant_count: number; non_compliant_count: number };
  critical_vulnerabilities: { total: number; affected_devices: number };
  open_incidents: { total: number; p1: number; p2: number; p3: number; p4: number };
  patch_compliance: { overall_rate: number };
  new_hires: { total: number; onboarded: number };
  terminations: { total: number; offboarded: number };
  sla_breaches: { active: number };
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  source: string;
}

const API_BASE = '/api/v1';

export async function fetchDevices(): Promise<PaginatedResponse<Device>> {
  const res = await fetch(`${API_BASE}/devices`);
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  const data = await res.json();
  // Handle both { devices: [...] } and { items: [...] } response formats
  if (data.devices && !data.items) {
    return { items: data.devices, total: data.total || data.devices.length, source: data.source || 'api' };
  }
  return data;
}

export async function fetchDevice(id: string): Promise<Device> {
  const res = await fetch(`${API_BASE}/devices/${encodeURIComponent(id)}`);
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  return res.json();
}

export async function fetchDashboardSummary(): Promise<DashboardSummary> {
  // For now, derive from devices endpoint. Later: dedicated dashboard API.
  const devices = await fetchDevices();
  const items = devices.items || [];
  const total = devices.total || items.length;
  const online = items.filter(d => d.status === 'active').length;
  const compliant = items.filter(d => d.compliance_status === 'compliant').length;
  return {
    total_devices: { value: total, previous_value: total - 5, change_percent: 0.4 },
    online_devices: { value: online, percentage: total > 0 ? Math.round((online / total) * 100) : 0 },
    offline_devices: { value: total - online, percentage: total > 0 ? Math.round(((total - online) / total) * 100) : 0 },
    compliance_rate: { rate: total > 0 ? Math.round((compliant / total) * 1000) / 10 : 0, compliant_count: compliant, non_compliant_count: total - compliant },
    critical_vulnerabilities: { total: 0, affected_devices: 0 },
    open_incidents: { total: 0, p1: 0, p2: 0, p3: 0, p4: 0 },
    patch_compliance: { overall_rate: 0 },
    new_hires: { total: 0, onboarded: 0 },
    terminations: { total: 0, offboarded: 0 },
    sla_breaches: { active: 0 },
  };
}

export async function checkHealth(): Promise<{ status: string; service: string }> {
  const res = await fetch('/health');
  if (!res.ok) throw new Error('API unreachable');
  return res.json();
}

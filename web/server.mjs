import http from 'http';
import { readFile } from 'fs/promises';
import { existsSync } from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const PORT = 3000;
const API_TARGET = 'http://localhost:8080';

const MIME_TYPES = {
  '.html': 'text/html',
  '.js': 'application/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.ico': 'image/x-icon',
};

const server = http.createServer(async (req, res) => {
  try {
    const url = req.url || '/';

    // Add CORS headers for all responses
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');

    // Handle preflight
    if (req.method === 'OPTIONS') {
      res.writeHead(204);
      res.end();
      return;
    }

    // Proxy API requests to backend
    if (url.startsWith('/api') || url.startsWith('/health')) {
      const targetUrl = API_TARGET + url;
      const proxyRes = await fetch(targetUrl, {
        method: req.method,
        headers: { ...req.headers },
      });
      const body = await proxyRes.text();
      res.writeHead(proxyRes.status, { 'Content-Type': 'application/json' });
      res.end(body);
      return;
    }

    // Serve static files from dist
    const filePath = path.join(__dirname, 'dist', url === '/' ? 'index.html' : url);
    const finalPath = existsSync(filePath) ? filePath : path.join(__dirname, 'dist', 'index.html');

    const ext = path.extname(finalPath);
    const contentType = MIME_TYPES[ext] || 'application/octet-stream';
    const content = await readFile(finalPath);

    res.writeHead(200, { 'Content-Type': contentType });
    res.end(content);
  } catch (err) {
    res.writeHead(500, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Internal Server Error', message: err.message }));
  }
});

server.listen(PORT, '0.0.0.0', () => {
  console.log(`Dashboard server running at http://0.0.0.0:${PORT}`);
  console.log(`API proxy: /api/* -> ${API_TARGET}`);
});

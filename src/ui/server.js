const http = require('http');
const fs = require('fs');
const path = require('path');

const PORT = process.env.PORT || 3000;
const DIST_DIR = path.join(__dirname, 'dist');

// MIME types
const mimeTypes = {
  '.html': 'text/html',
  '.js': 'application/javascript',
  '.css': 'text/css',
  '.json': 'application/json',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.gif': 'image/gif',
  '.svg': 'image/svg+xml',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.eot': 'application/vnd.ms-fontobject',
};

function getMimeType(filePath) {
  const ext = path.extname(filePath).toLowerCase();
  return mimeTypes[ext] || 'application/octet-stream';
}

function serveFile(filePath, res) {
  fs.readFile(filePath, (err, data) => {
    if (err) {
      res.writeHead(404, { 'Content-Type': 'text/plain' });
      res.end('Not Found');
      return;
    }

    const mimeType = getMimeType(filePath);
    const headers = { 'Content-Type': mimeType };
    
    // Add CORS headers for font files to prevent cross-origin issues
    if (/\.(woff|woff2|ttf|eot|otf)$/i.test(filePath)) {
      headers['Access-Control-Allow-Origin'] = '*';
      headers['Cache-Control'] = 'public, max-age=31536000';
    }
    
    res.writeHead(200, headers);
    res.end(data);
  });
}

const server = http.createServer((req, res) => {
  // Health check endpoint
  if (req.url === '/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({
      status: 'ok',
      timestamp: new Date().toISOString(),
      service: process.env.VITE_SERVICE_NAME,
    }));
    return;
  }

  // Determine file path
  let filePath = path.join(DIST_DIR, req.url === '/' ? 'index.html' : req.url);

  // Security: prevent directory traversal
  if (!filePath.startsWith(DIST_DIR)) {
    res.writeHead(403, { 'Content-Type': 'text/plain' });
    res.end('Forbidden');
    return;
  }

  // Check if file exists
  fs.stat(filePath, (err, stats) => {
    if (err || !stats.isFile()) {
      // If file doesn't exist, serve index.html for SPA routing
      filePath = path.join(DIST_DIR, 'index.html');
    }

    serveFile(filePath, res);
  });
});

server.listen(PORT, '0.0.0.0', () => {
  console.log(`UI server running on http://0.0.0.0:${PORT}`);
  console.log(`Serving static files from: ${DIST_DIR}`);
});

const API_BASE_URL = window.location.origin;

// Helper function for API requests
async function apiRequest(url, options = {}) {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.status} ${response.statusText}`);
  }

  return response.json();
}

// Fetch all clients
export const fetchClients = () => {
  return apiRequest('/api/clients');
};

// Fetch a specific client
export const fetchClient = (clientId) => {
  return apiRequest(`/api/clients/${clientId}`);
};

// Fetch network requests for a client
export const fetchClientRequests = (clientId, limit = 100) => {
  return apiRequest(`/api/clients/${clientId}/requests?limit=${limit}`);
};

// Fetch server configuration
export const fetchConfig = () => {
  return apiRequest('/api/config');
};

// Update server configuration
export const updateConfig = (config) => {
  return apiRequest('/api/config', {
    method: 'PUT',
    body: JSON.stringify(config),
  });
};

// Update client configuration
export const updateClientConfig = (clientId, config) => {
  return apiRequest(`/api/clients/${clientId}/config`, {
    method: 'PUT',
    body: JSON.stringify(config),
  });
};

// Send a command to a client
export const sendClientCommand = (clientId, command) => {
  return apiRequest(`/api/clients/${clientId}/command`, {
    method: 'POST',
    body: JSON.stringify(command),
  });
};

// Get a client configuration file
export const getClientConfigFile = (clientId, filePath) => {
  return apiRequest(`/api/clients/${clientId}/files`, {
    method: 'POST',
    body: JSON.stringify({ path: filePath }),
  });
};

// Update a client configuration file
export const updateClientConfigFile = (clientId, filePath, content) => {
  return apiRequest(`/api/clients/${clientId}/files`, {
    method: 'PUT',
    body: JSON.stringify({ path: filePath, content }),
  });
};
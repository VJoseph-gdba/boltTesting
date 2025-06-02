import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Box from '@mui/material/Box';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Chip from '@mui/material/Chip';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import { LineChart } from '@mui/x-charts/LineChart';
import { fetchClient, fetchClientRequests, getClientConfigFile, updateClientConfigFile } from '../services/api';
import { formatDistanceToNow, format } from 'date-fns';
import { 
  Laptop, 
  Clock, 
  Activity, 
  CheckCircle2, 
  XCircle,
  File,
  FileEdit,
  Save
} from 'lucide-react';

function ClientDetails() {
  const { id } = useParams();
  const [client, setClient] = useState(null);
  const [requests, setRequests] = useState([]);
  const [tabValue, setTabValue] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [configFileOpen, setConfigFileOpen] = useState(false);
  const [configFile, setConfigFile] = useState({ name: '', content: '' });
  const [chartData, setChartData] = useState({
    timeLabels: [],
    responseData: [],
    dnsData: [],
    tcpData: [],
    tlsData: []
  });

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        // Fetch client info
        const clientData = await fetchClient(id);
        setClient(clientData);
        
        // Fetch client requests
        const requestsData = await fetchClientRequests(id, 100);
        setRequests(requestsData);
        
        // Prepare chart data
        if (requestsData.length > 0) {
          // Sort by time
          requestsData.sort((a, b) => new Date(a.startTime) - new Date(b.startTime));
          
          // Take last 20 entries
          const chartEntries = requestsData.slice(-20);
          
          setChartData({
            timeLabels: chartEntries.map(r => format(new Date(r.startTime), 'HH:mm:ss')),
            responseData: chartEntries.map(r => r.error ? null : r.totalTime),
            dnsData: chartEntries.map(r => r.error ? null : r.dnsTime),
            tcpData: chartEntries.map(r => r.error ? null : r.tcpTime),
            tlsData: chartEntries.map(r => r.error ? null : r.tlsTime)
          });
        }
      } catch (error) {
        console.error('Error fetching client data:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 30000); // Refresh every 30 seconds
    
    return () => clearInterval(interval);
  }, [id]);

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleOpenConfigFile = async () => {
    try {
      // This is a placeholder - in a real app, we would get the actual file
      const configData = await getClientConfigFile(id, 'client.json');
      setConfigFile(configData);
      setConfigFileOpen(true);
    } catch (error) {
      console.error('Error fetching config file:', error);
      // Fallback for demo
      setConfigFile({
        name: 'client.json',
        content: JSON.stringify({
          serverAddress: "http://localhost:8080",
          clientName: "Demo Client",
          targets: [
            {
              name: "Google",
              url: "https://www.google.com",
              interval: 60,
              enabled: true
            },
            {
              name: "GitHub",
              url: "https://github.com",
              interval: 60,
              enabled: true
            }
          ],
          logLevel: "info"
        }, null, 2)
      });
      setConfigFileOpen(true);
    }
  };

  const handleCloseConfigFile = () => {
    setConfigFileOpen(false);
  };

  const handleSaveConfigFile = async () => {
    try {
      await updateClientConfigFile(id, configFile.name, configFile.content);
      handleCloseConfigFile();
    } catch (error) {
      console.error('Error saving config file:', error);
      // For demo, just close the dialog
      handleCloseConfigFile();
    }
  };

  const handleConfigContentChange = (event) => {
    setConfigFile({ ...configFile, content: event.target.value });
  };

  // Calculate statistics
  const stats = {
    totalRequests: requests.length,
    successRequests: requests.filter(r => !r.error && r.statusCode < 400).length,
    errorRequests: requests.filter(r => r.error || r.statusCode >= 400).length,
    avgResponseTime: requests.length > 0 
      ? Math.round(requests.filter(r => !r.error).reduce((sum, r) => sum + r.totalTime, 0) / 
                  requests.filter(r => !r.error).length) 
      : 0,
    avgDnsTime: requests.length > 0 
      ? Math.round(requests.filter(r => !r.error).reduce((sum, r) => sum + r.dnsTime, 0) / 
                  requests.filter(r => !r.error).length) 
      : 0,
    avgTcpTime: requests.length > 0 
      ? Math.round(requests.filter(r => !r.error).reduce((sum, r) => sum + r.tcpTime, 0) / 
                  requests.filter(r => !r.error).length) 
      : 0,
    avgTlsTime: requests.length > 0 
      ? Math.round(requests.filter(r => !r.error).reduce((sum, r) => sum + r.tlsTime, 0) / 
                  requests.filter(r => !r.error).length) 
      : 0
  };

  if (isLoading || !client) {
    return (
      <Box sx={{ p: 3 }}>
        <Typography>Loading client data...</Typography>
      </Box>
    );
  }

  return (
    <Box>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Client Details
        </Typography>
        
        <Paper 
          elevation={2} 
          sx={{ 
            p: 3, 
            borderRadius: 2,
            background: theme => 
              client.status === 'online' 
                ? `linear-gradient(90deg, ${theme.palette.background.paper} 0%, ${theme.palette.success.dark}22 100%)` 
                : `linear-gradient(90deg, ${theme.palette.background.paper} 0%, ${theme.palette.error.dark}22 100%)`
          }}
        >
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Box display="flex" alignItems="center">
                <Laptop size={36} style={{ marginRight: '16px' }} />
                <Box>
                  <Typography variant="h5" component="h2">
                    {client.name}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {client.ipAddress} • {client.osInfo}
                  </Typography>
                </Box>
              </Box>
              
              <Box display="flex" alignItems="center" mt={2}>
                <Chip 
                  label={client.status === 'online' ? 'Online' : 'Offline'} 
                  color={client.status === 'online' ? 'success' : 'error'}
                  sx={{ mr: 1 }}
                />
                <Typography variant="body2" color="text.secondary">
                  <Clock size={16} style={{ verticalAlign: 'text-bottom', marginRight: '4px' }} />
                  Last seen {formatDistanceToNow(new Date(client.lastSeen), { addSuffix: true })}
                </Typography>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Connected Since
                  </Typography>
                  <Typography variant="body1">
                    {client.status === 'online' 
                      ? formatDistanceToNow(new Date(client.connectedAt)) 
                      : 'Offline'}
                  </Typography>
                </Grid>
                
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Version
                  </Typography>
                  <Typography variant="body1">
                    {client.version}
                  </Typography>
                </Grid>
                
                <Grid item xs={12}>
                  <Box display="flex" justifyContent="flex-end">
                    <Button 
                      variant="outlined" 
                      startIcon={<FileEdit size={16} />}
                      onClick={handleOpenConfigFile}
                    >
                      Edit Configuration
                    </Button>
                  </Box>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Paper>
      </Box>
      
      <Tabs 
        value={tabValue} 
        onChange={handleTabChange}
        sx={{ mb: 3, borderBottom: 1, borderColor: 'divider' }}
      >
        <Tab label="Overview" />
        <Tab label="Network Requests" />
        <Tab label="Configuration" />
      </Tabs>
      
      {/* Overview Tab */}
      {tabValue === 0 && (
        <Box>
          <Grid container spacing={3} sx={{ mb: 4 }}>
            <Grid item xs={12} sm={6} md={3}>
              <Card 
                elevation={2}
                sx={{ 
                  borderRadius: 2,
                  height: '100%',
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    <Activity size={20} color="#2196F3" />
                    <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                      Requests
                    </Typography>
                  </Box>
                  <Typography variant="h4" component="div">
                    {stats.totalRequests}
                  </Typography>
                  <Box display="flex" mt={1}>
                    <Chip 
                      label={`${stats.successRequests} Success`} 
                      size="small" 
                      color="success" 
                      sx={{ mr: 1 }} 
                    />
                    <Chip 
                      label={`${stats.errorRequests} Error`} 
                      size="small" 
                      color="error" 
                    />
                  </Box>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card 
                elevation={2}
                sx={{ 
                  borderRadius: 2,
                  height: '100%',
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    <Clock size={20} color="#00BCD4" />
                    <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                      Avg Response
                    </Typography>
                  </Box>
                  <Typography variant="h4" component="div">
                    {stats.avgResponseTime}ms
                  </Typography>
                  <Typography variant="body2" color="text.secondary" mt={1}>
                    From successful requests
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card 
                elevation={2}
                sx={{ 
                  borderRadius: 2,
                  height: '100%',
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    <CheckCircle2 size={20} color="#4CAF50" />
                    <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                      Success Rate
                    </Typography>
                  </Box>
                  <Typography variant="h4" component="div">
                    {stats.totalRequests > 0 
                      ? Math.round((stats.successRequests / stats.totalRequests) * 100) 
                      : 0}%
                  </Typography>
                  <Typography variant="body2" color="text.secondary" mt={1}>
                    Overall success rate
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
            
            <Grid item xs={12} sm={6} md={3}>
              <Card 
                elevation={2}
                sx={{ 
                  borderRadius: 2,
                  height: '100%',
                }}
              >
                <CardContent>
                  <Box display="flex" alignItems="center" mb={1}>
                    <XCircle size={20} color="#F44336" />
                    <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                      DNS Time
                    </Typography>
                  </Box>
                  <Typography variant="h4" component="div">
                    {stats.avgDnsTime}ms
                  </Typography>
                  <Typography variant="body2" color="text.secondary" mt={1}>
                    Average DNS resolution time
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          </Grid>
          
          <Paper 
            elevation={2} 
            sx={{ 
              p: 3, 
              borderRadius: 2,
              mb: 4,
              height: 400,
            }}
          >
            <Typography variant="h6" component="h3" gutterBottom>
              Response Time Trends
            </Typography>
            
            {chartData.timeLabels.length > 0 ? (
              <LineChart
                series={[
                  {
                    data: chartData.responseData,
                    label: 'Total Response Time',
                    color: '#2196F3',
                  },
                  {
                    data: chartData.dnsData,
                    label: 'DNS Time',
                    color: '#FFC107',
                  },
                  {
                    data: chartData.tcpData,
                    label: 'TCP Time',
                    color: '#4CAF50',
                  },
                  {
                    data: chartData.tlsData,
                    label: 'TLS Time',
                    color: '#9C27B0',
                  },
                ]}
                xAxis={[{ 
                  scaleType: 'point',
                  data: chartData.timeLabels,
                  label: 'Time',
                }]}
                yAxis={[{
                  label: 'Time (ms)',
                }]}
                height={300}
                margin={{ top: 20, right: 40, bottom: 50, left: 40 }}
                sx={{
                  '.MuiLineElement-root': {
                    strokeWidth: 2,
                  },
                  '.MuiMarkElement-root': {
                    stroke: 'white',
                    strokeWidth: 2,
                    fill: 'currentColor',
                    r: 4,
                  },
                }}
              />
            ) : (
              <Box display="flex" justifyContent="center" alignItems="center" height="100%">
                <Typography color="text.secondary">
                  No data available for chart
                </Typography>
              </Box>
            )}
          </Paper>
          
          <Paper 
            elevation={2} 
            sx={{ 
              p: 3, 
              borderRadius: 2,
            }}
          >
            <Typography variant="h6" component="h3" gutterBottom>
              Recent Network Activity
            </Typography>
            
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Target</TableCell>
                    <TableCell>URL</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Time</TableCell>
                    <TableCell>Duration</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {requests.slice(0, 5).map((request) => (
                    <TableRow key={request.id} hover>
                      <TableCell>{request.targetName}</TableCell>
                      <TableCell sx={{ maxWidth: 300, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                        {request.url}
                      </TableCell>
                      <TableCell>
                        {request.error ? (
                          <Chip 
                            icon={<XCircle size={14} />}
                            label={request.errorType || 'Error'} 
                            size="small"
                            color="error"
                          />
                        ) : (
                          <Chip 
                            icon={<CheckCircle2 size={14} />}
                            label={request.statusCode} 
                            size="small"
                            color={request.statusCode < 400 ? 'success' : 'error'}
                          />
                        )}
                      </TableCell>
                      <TableCell>
                        {formatDistanceToNow(new Date(request.startTime), { addSuffix: true })}
                      </TableCell>
                      <TableCell>
                        {request.totalTime}ms
                      </TableCell>
                    </TableRow>
                  ))}
                  
                  {requests.length === 0 && (
                    <TableRow>
                      <TableCell colSpan={5} align="center">
                        <Typography color="text.secondary">
                          No requests found
                        </Typography>
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </TableContainer>
          </Paper>
        </Box>
      )}
      
      {/* Network Requests Tab */}
      {tabValue === 1 && (
        <Paper 
          elevation={2} 
          sx={{ 
            p: 3, 
            borderRadius: 2,
          }}
        >
          <Typography variant="h6" component="h3" gutterBottom>
            Network Requests
          </Typography>
          
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Target</TableCell>
                  <TableCell>URL</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Time</TableCell>
                  <TableCell>DNS</TableCell>
                  <TableCell>TCP</TableCell>
                  <TableCell>TLS</TableCell>
                  <TableCell>Total</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {requests.map((request) => (
                  <TableRow key={request.id} hover>
                    <TableCell>{request.targetName}</TableCell>
                    <TableCell sx={{ maxWidth: 250, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                      {request.url}
                    </TableCell>
                    <TableCell>
                      {request.error ? (
                        <Chip 
                          icon={<XCircle size={14} />}
                          label={request.errorType || 'Error'} 
                          size="small"
                          color="error"
                        />
                      ) : (
                        <Chip 
                          icon={<CheckCircle2 size={14} />}
                          label={request.statusCode} 
                          size="small"
                          color={request.statusCode < 400 ? 'success' : 'error'}
                        />
                      )}
                    </TableCell>
                    <TableCell>
                      {format(new Date(request.startTime), 'yyyy-MM-dd HH:mm:ss')}
                    </TableCell>
                    <TableCell>{request.dnsTime}ms</TableCell>
                    <TableCell>{request.tcpTime}ms</TableCell>
                    <TableCell>{request.tlsTime}ms</TableCell>
                    <TableCell>{request.totalTime}ms</TableCell>
                  </TableRow>
                ))}
                
                {requests.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={8} align="center">
                      <Typography color="text.secondary">
                        No requests found
                      </Typography>
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      )}
      
      {/* Configuration Tab */}
      {tabValue === 2 && (
        <Paper 
          elevation={2} 
          sx={{ 
            p: 3, 
            borderRadius: 2,
          }}
        >
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
            <Typography variant="h6" component="h3">
              Configuration Files
            </Typography>
            <Button 
              variant="contained" 
              startIcon={<FileEdit size={16} />}
              onClick={handleOpenConfigFile}
            >
              Edit Configuration
            </Button>
          </Box>
          
          <TableContainer>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>File</TableCell>
                  <TableCell>Path</TableCell>
                  <TableCell>Last Modified</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                <TableRow hover>
                  <TableCell>
                    <Box display="flex" alignItems="center">
                      <File size={20} style={{ marginRight: '8px' }} />
                      client.json
                    </Box>
                  </TableCell>
                  <TableCell>~/.config/NetworkMonitor/client.json</TableCell>
                  <TableCell>
                    {formatDistanceToNow(new Date(), { addSuffix: true })}
                  </TableCell>
                  <TableCell>
                    <Button 
                      variant="outlined" 
                      size="small"
                      onClick={handleOpenConfigFile}
                    >
                      Edit
                    </Button>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </Paper>
      )}
      
      {/* Config File Edit Dialog */}
      <Dialog
        open={configFileOpen}
        onClose={handleCloseConfigFile}
        fullWidth
        maxWidth="md"
      >
        <DialogTitle>
          <Box display="flex" alignItems="center">
            <File size={20} style={{ marginRight: '8px' }} />
            Edit Configuration: {configFile.name}
          </Box>
        </DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            fullWidth
            multiline
            rows={20}
            variant="outlined"
            value={configFile.content}
            onChange={handleConfigContentChange}
            InputProps={{
              style: { fontFamily: 'monospace' }
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseConfigFile}>Cancel</Button>
          <Button 
            variant="contained" 
            onClick={handleSaveConfigFile}
            startIcon={<Save size={16} />}
          >
            Save
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default ClientDetails;
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import Typography from '@mui/material/Typography';
import Box from '@mui/material/Box';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Button from '@mui/material/Button';
import Chip from '@mui/material/Chip';
import { 
  Activity,
  Laptop,
  AlertCircle,
  Clock,
  CheckCircle2,
  XCircle,
  ArrowRight
} from 'lucide-react';
import { fetchClients, fetchClientRequests } from '../services/api';
import { formatDistanceToNow } from 'date-fns';

function Dashboard() {
  const navigate = useNavigate();
  const [clients, setClients] = useState([]);
  const [recentRequests, setRecentRequests] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [stats, setStats] = useState({
    totalClients: 0,
    onlineClients: 0,
    offlineClients: 0,
    successRate: 0,
    totalRequests: 0,
    errorRequests: 0
  });

  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      try {
        // Fetch clients
        const clientsData = await fetchClients();
        setClients(clientsData);
        
        // Calculate stats
        const onlineCount = clientsData.filter(c => c.status === 'online').length;
        
        // Fetch recent requests for online clients
        let allRequests = [];
        for (const client of clientsData.filter(c => c.status === 'online').slice(0, 3)) {
          const requests = await fetchClientRequests(client.id, 10);
          allRequests = [...allRequests, ...requests.map(r => ({ ...r, clientName: client.name, clientId: client.id }))];
        }
        
        // Sort by time
        allRequests.sort((a, b) => new Date(b.startTime) - new Date(a.startTime));
        setRecentRequests(allRequests.slice(0, 10));
        
        // Calculate success rate
        const totalReqs = allRequests.length;
        const errorReqs = allRequests.filter(r => r.error).length;
        const successRate = totalReqs > 0 ? Math.round(((totalReqs - errorReqs) / totalReqs) * 100) : 100;
        
        setStats({
          totalClients: clientsData.length,
          onlineClients: onlineCount,
          offlineClients: clientsData.length - onlineCount,
          successRate: successRate,
          totalRequests: totalReqs,
          errorRequests: errorReqs
        });
        
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
    const interval = setInterval(fetchData, 30000); // Refresh every 30 seconds
    
    return () => clearInterval(interval);
  }, []);

  return (
    <Box>
      <Typography variant="h4" component="h1" gutterBottom>
        Dashboard
      </Typography>
      
      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card 
            elevation={2}
            sx={{ 
              borderRadius: 2,
              height: '100%',
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-4px)',
              }
            }}
          >
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <Laptop size={24} color="#2196F3" />
                <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                  Total Clients
                </Typography>
              </Box>
              <Typography variant="h3" component="div">
                {stats.totalClients}
              </Typography>
              <Box display="flex" mt={1}>
                <Chip 
                  label={`${stats.onlineClients} Online`} 
                  size="small" 
                  color="success" 
                  sx={{ mr: 1 }} 
                />
                <Chip 
                  label={`${stats.offlineClients} Offline`} 
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
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-4px)',
              }
            }}
          >
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <Activity size={24} color="#00BCD4" />
                <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                  Success Rate
                </Typography>
              </Box>
              <Typography variant="h3" component="div">
                {stats.successRate}%
              </Typography>
              <Box display="flex" mt={1}>
                <Typography variant="body2" color="text.secondary">
                  From {stats.totalRequests} requests
                </Typography>
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
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-4px)',
              }
            }}
          >
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <CheckCircle2 size={24} color="#4CAF50" />
                <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                  Healthy Clients
                </Typography>
              </Box>
              <Typography variant="h3" component="div">
                {stats.onlineClients}
              </Typography>
              <Box display="flex" mt={1}>
                <Typography variant="body2" color="text.secondary">
                  {Math.round((stats.onlineClients / stats.totalClients) * 100) || 0}% of total
                </Typography>
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
              transition: 'transform 0.3s',
              '&:hover': {
                transform: 'translateY(-4px)',
              }
            }}
          >
            <CardContent>
              <Box display="flex" alignItems="center" mb={2}>
                <AlertCircle size={24} color="#F44336" />
                <Typography variant="h6" component="div" sx={{ ml: 1 }}>
                  Errors
                </Typography>
              </Box>
              <Typography variant="h3" component="div">
                {stats.errorRequests}
              </Typography>
              <Box display="flex" mt={1}>
                <Typography variant="body2" color="text.secondary">
                  {Math.round((stats.errorRequests / stats.totalRequests) * 100) || 0}% of requests
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      {/* Client List */}
      <Paper 
        elevation={2} 
        sx={{ 
          p: 3, 
          mb: 4, 
          borderRadius: 2,
          transition: 'box-shadow 0.3s',
          '&:hover': {
            boxShadow: 6,
          }
        }}
      >
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5" component="h2">
            Active Clients
          </Typography>
          <Button 
            variant="outlined" 
            color="primary" 
            endIcon={<ArrowRight size={16} />}
          >
            View All
          </Button>
        </Box>
        
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Client</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>IP Address</TableCell>
                <TableCell>Last Seen</TableCell>
                <TableCell>Uptime</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {clients
                .filter(client => client.status === 'online')
                .slice(0, 5)
                .map((client) => (
                  <TableRow key={client.id} hover>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Laptop size={20} style={{ marginRight: '8px' }} />
                        {client.name}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip 
                        label={client.status === 'online' ? 'Online' : 'Offline'} 
                        size="small"
                        color={client.status === 'online' ? 'success' : 'error'}
                      />
                    </TableCell>
                    <TableCell>{client.ipAddress}</TableCell>
                    <TableCell>
                      <Box display="flex" alignItems="center">
                        <Clock size={16} style={{ marginRight: '4px' }} />
                        {formatDistanceToNow(new Date(client.lastSeen), { addSuffix: true })}
                      </Box>
                    </TableCell>
                    <TableCell>
                      {client.status === 'online' ? 
                        formatDistanceToNow(new Date(client.connectedAt)) : 'N/A'}
                    </TableCell>
                    <TableCell>
                      <Button 
                        variant="contained" 
                        size="small"
                        onClick={() => navigate(`/clients/${client.id}`)}
                      >
                        View
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
                
              {clients.filter(client => client.status === 'online').length === 0 && (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    <Typography color="text.secondary">
                      No online clients found
                    </Typography>
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>
      
      {/* Recent Requests */}
      <Paper 
        elevation={2} 
        sx={{ 
          p: 3, 
          borderRadius: 2,
          transition: 'box-shadow 0.3s',
          '&:hover': {
            boxShadow: 6,
          }
        }}
      >
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5" component="h2">
            Recent Network Requests
          </Typography>
          <Button 
            variant="outlined" 
            color="primary" 
            endIcon={<ArrowRight size={16} />}
          >
            View All
          </Button>
        </Box>
        
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>URL</TableCell>
                <TableCell>Client</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Time</TableCell>
                <TableCell>Duration</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {recentRequests.slice(0, 5).map((request) => (
                <TableRow key={request.id} hover>
                  <TableCell sx={{ maxWidth: 300, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {request.url}
                  </TableCell>
                  <TableCell>
                    <Box 
                      sx={{ 
                        cursor: 'pointer',
                        '&:hover': { textDecoration: 'underline' }
                      }}
                      onClick={() => navigate(`/clients/${request.clientId}`)}
                    >
                      {request.clientName || 'Unknown'}
                    </Box>
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
              
              {recentRequests.length === 0 && (
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
  );
}

export default Dashboard;
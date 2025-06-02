import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import Box from '@mui/material/Box';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Divider from '@mui/material/Divider';
import Typography from '@mui/material/Typography';
import Badge from '@mui/material/Badge';
import { 
  LayoutDashboard, 
  Users, 
  Settings, 
  History, 
  Activity,
  Laptop,
  Server,
  ChevronRight
} from 'lucide-react';
import { fetchClients } from '../services/api';

const drawerWidth = 240;

function Sidebar({ open, onClose }) {
  const navigate = useNavigate();
  const location = useLocation();
  const [clients, setClients] = useState([]);

  useEffect(() => {
    const getClients = async () => {
      try {
        const data = await fetchClients();
        setClients(data);
      } catch (error) {
        console.error('Failed to fetch clients:', error);
      }
    };

    getClients();
    const interval = setInterval(getClients, 30000); // Refresh every 30 seconds
    
    return () => clearInterval(interval);
  }, []);

  const onlineCount = clients.filter(client => client.status === 'online').length;

  const navigateTo = (path) => {
    navigate(path);
    if (window.innerWidth < 600) {
      onClose();
    }
  };

  const isActive = (path) => location.pathname === path;

  return (
    <Drawer
      variant="persistent"
      open={open}
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: drawerWidth,
          boxSizing: 'border-box',
          top: '64px',
          height: 'calc(100% - 64px)',
        },
      }}
    >
      <Box sx={{ overflow: 'auto', height: '100%', display: 'flex', flexDirection: 'column' }}>
        <List>
          <ListItem disablePadding>
            <ListItemButton 
              selected={isActive('/')}
              onClick={() => navigateTo('/')}
            >
              <ListItemIcon>
                <LayoutDashboard size={24} />
              </ListItemIcon>
              <ListItemText primary="Dashboard" />
            </ListItemButton>
          </ListItem>
          
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <Badge badgeContent={onlineCount} color="primary">
                  <Users size={24} />
                </Badge>
              </ListItemIcon>
              <ListItemText primary="Clients" />
            </ListItemButton>
          </ListItem>
          
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <Activity size={24} />
              </ListItemIcon>
              <ListItemText primary="Network Status" />
            </ListItemButton>
          </ListItem>
          
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <History size={24} />
              </ListItemIcon>
              <ListItemText primary="History" />
            </ListItemButton>
          </ListItem>
          
          <ListItem disablePadding>
            <ListItemButton>
              <ListItemIcon>
                <Settings size={24} />
              </ListItemIcon>
              <ListItemText primary="Settings" />
            </ListItemButton>
          </ListItem>
        </List>
        
        <Divider />
        
        <Box sx={{ px: 2, py: 1 }}>
          <Typography variant="subtitle2" color="text.secondary">
            CONNECTED CLIENTS
          </Typography>
        </Box>
        
        <List sx={{ overflow: 'auto', flexGrow: 1 }}>
          {clients
            .filter(client => client.status === 'online')
            .map((client) => (
              <ListItem key={client.id} disablePadding>
                <ListItemButton 
                  selected={isActive(`/clients/${client.id}`)}
                  onClick={() => navigateTo(`/clients/${client.id}`)}
                >
                  <ListItemIcon>
                    <Laptop size={20} />
                  </ListItemIcon>
                  <ListItemText 
                    primary={client.name} 
                    primaryTypographyProps={{ noWrap: true }}
                  />
                  <ChevronRight size={16} />
                </ListItemButton>
              </ListItem>
            ))}
        </List>
        
        <Divider />
        
        <Box sx={{ px: 2, py: 1 }}>
          <Typography variant="subtitle2" color="text.secondary">
            DISCONNECTED CLIENTS
          </Typography>
        </Box>
        
        <List sx={{ overflow: 'auto', flexGrow: 1 }}>
          {clients
            .filter(client => client.status !== 'online')
            .map((client) => (
              <ListItem key={client.id} disablePadding>
                <ListItemButton 
                  selected={isActive(`/clients/${client.id}`)}
                  onClick={() => navigateTo(`/clients/${client.id}`)}
                  sx={{ opacity: 0.6 }}
                >
                  <ListItemIcon>
                    <Laptop size={20} />
                  </ListItemIcon>
                  <ListItemText 
                    primary={client.name} 
                    primaryTypographyProps={{ noWrap: true }}
                  />
                  <ChevronRight size={16} />
                </ListItemButton>
              </ListItem>
            ))}
        </List>
        
        <Divider />
        
        <Box sx={{ p: 2, display: 'flex', alignItems: 'center' }}>
          <Server size={20} />
          <Typography variant="body2" sx={{ ml: 1 }}>
            Server Status: Online
          </Typography>
        </Box>
      </Box>
    </Drawer>
  );
}

export default Sidebar;
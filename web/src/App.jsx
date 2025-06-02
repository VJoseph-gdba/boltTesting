import React, { useState, useEffect } from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Dashboard from './pages/Dashboard';
import ClientDetails from './pages/ClientDetails';
import Sidebar from './components/Sidebar';
import TopBar from './components/TopBar';
import Box from '@mui/material/Box';
import useMediaQuery from '@mui/material/useMediaQuery';

const lightTheme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: '#2196F3',
    },
    secondary: {
      main: '#00BCD4',
    },
    success: {
      main: '#4CAF50',
    },
    warning: {
      main: '#FF9800',
    },
    error: {
      main: '#F44336',
    },
    background: {
      default: '#f5f5f5',
      paper: '#ffffff',
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
    h1: {
      fontWeight: 500,
      fontSize: '2.5rem',
      lineHeight: 1.2,
    },
    h2: {
      fontWeight: 500,
      fontSize: '2rem',
      lineHeight: 1.2,
    },
    h3: {
      fontWeight: 500,
      fontSize: '1.75rem',
      lineHeight: 1.2,
    },
    h4: {
      fontWeight: 500,
      fontSize: '1.5rem',
      lineHeight: 1.2,
    },
    h5: {
      fontWeight: 500,
      fontSize: '1.25rem',
      lineHeight: 1.2,
    },
    h6: {
      fontWeight: 500,
      fontSize: '1rem',
      lineHeight: 1.2,
    },
    body1: {
      lineHeight: 1.5,
    },
    body2: {
      lineHeight: 1.5,
    },
  },
  shape: {
    borderRadius: 8,
  },
});

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90CAF9',
    },
    secondary: {
      main: '#80DEEA',
    },
    success: {
      main: '#81C784',
    },
    warning: {
      main: '#FFB74D',
    },
    error: {
      main: '#E57373',
    },
    background: {
      default: '#121212',
      paper: '#1E1E1E',
    },
  },
  typography: lightTheme.typography,
  shape: lightTheme.shape,
});

function App() {
  const [darkMode, setDarkMode] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const prefersDarkMode = useMediaQuery('(prefers-color-scheme: dark)');

  useEffect(() => {
    const savedMode = localStorage.getItem('darkMode');
    if (savedMode !== null) {
      setDarkMode(savedMode === 'true');
    } else {
      setDarkMode(prefersDarkMode);
    }
  }, [prefersDarkMode]);

  const toggleDarkMode = () => {
    const newMode = !darkMode;
    setDarkMode(newMode);
    localStorage.setItem('darkMode', newMode.toString());
  };

  const toggleDrawer = () => {
    setDrawerOpen(!drawerOpen);
  };

  return (
    <ThemeProvider theme={darkMode ? darkTheme : lightTheme}>
      <CssBaseline />
      <BrowserRouter>
        <Box sx={{ display: 'flex', height: '100vh' }}>
          <TopBar
            darkMode={darkMode}
            toggleDarkMode={toggleDarkMode}
            toggleDrawer={toggleDrawer}
          />
          <Sidebar open={drawerOpen} onClose={toggleDrawer} />
          <Box
            component="main"
            sx={{
              flexGrow: 1,
              p: 3,
              pt: 10,
              width: { xs: '100%', sm: `calc(100% - ${drawerOpen ? 240 : 0}px)` },
              marginLeft: { xs: 0, sm: drawerOpen ? 0 : -240 },
              transition: 'margin 225ms cubic-bezier(0.4, 0, 0.6, 1) 0ms',
            }}
          >
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/clients/:id" element={<ClientDetails />} />
            </Routes>
          </Box>
        </Box>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
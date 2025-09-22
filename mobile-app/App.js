import React, { useState, useEffect } from 'react';
import {
  SafeAreaView,
  ScrollView,
  StatusBar,
  StyleSheet,
  Text,
  View,
  TouchableOpacity,
  Alert,
  RefreshControl,
} from 'react-native';
import AsyncStorage from '@react-native-async-storage/async-storage';
import PushNotification from 'react-native-push-notification';
import axios from 'axios';

const App = () => {
  const [processes, setProcesses] = useState([]);
  const [serverEndpoint, setServerEndpoint] = useState('http://localhost:8080');
  const [isConnected, setIsConnected] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  useEffect(() => {
    initializeApp();
    setupNotifications();
  }, []);

  const initializeApp = async () => {
    try {
      const savedEndpoint = await AsyncStorage.getItem('serverEndpoint');
      if (savedEndpoint) {
        setServerEndpoint(savedEndpoint);
      }
      await fetchProcesses();
    } catch (error) {
      console.error('App initialization error:', error);
    }
  };

  const setupNotifications = () => {
    PushNotification.configure({
      onNotification: function(notification) {
        console.log('Notification received:', notification);
      },
      requestPermissions: Platform.OS === 'ios',
    });

    PushNotification.createChannel(
      {
        channelId: 'gproc-alerts',
        channelName: 'GProc Alerts',
        channelDescription: 'Process monitoring alerts',
        soundName: 'default',
        importance: 4,
        vibrate: true,
      },
      (created) => console.log('Notification channel created:', created)
    );
  };

  const fetchProcesses = async () => {
    try {
      const response = await axios.get(`${serverEndpoint}/api/v1/processes`, {
        timeout: 5000,
      });
      setProcesses(response.data.processes || []);
      setIsConnected(true);
    } catch (error) {
      console.error('Failed to fetch processes:', error);
      setIsConnected(false);
      showAlert('Connection Error', 'Failed to connect to GProc server');
    }
  };

  const onRefresh = async () => {
    setRefreshing(true);
    await fetchProcesses();
    setRefreshing(false);
  };

  const startProcess = async (processName) => {
    try {
      await axios.post(`${serverEndpoint}/api/v1/processes/${processName}/start`);
      await fetchProcesses();
      showNotification('Process Started', `${processName} has been started`);
    } catch (error) {
      showAlert('Error', `Failed to start ${processName}`);
    }
  };

  const stopProcess = async (processName) => {
    try {
      await axios.post(`${serverEndpoint}/api/v1/processes/${processName}/stop`);
      await fetchProcesses();
      showNotification('Process Stopped', `${processName} has been stopped`);
    } catch (error) {
      showAlert('Error', `Failed to stop ${processName}`);
    }
  };

  const showAlert = (title, message) => {
    Alert.alert(title, message);
  };

  const showNotification = (title, message) => {
    PushNotification.localNotification({
      channelId: 'gproc-alerts',
      title: title,
      message: message,
      playSound: true,
      soundName: 'default',
    });
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'running':
        return '#4CAF50';
      case 'stopped':
        return '#9E9E9E';
      case 'failed':
        return '#F44336';
      default:
        return '#FF9800';
    }
  };

  const ProcessCard = ({ process }) => (
    <View style={styles.processCard}>
      <View style={styles.processHeader}>
        <Text style={styles.processName}>{process.name}</Text>
        <View style={[styles.statusBadge, { backgroundColor: getStatusColor(process.status) }]}>
          <Text style={styles.statusText}>{process.status}</Text>
        </View>
      </View>
      
      <Text style={styles.processCommand}>{process.command}</Text>
      
      <View style={styles.processMetrics}>
        <Text style={styles.metricText}>CPU: {process.cpu_usage || 0}%</Text>
        <Text style={styles.metricText}>Memory: {Math.round((process.memory_usage || 0) / 1024 / 1024)}MB</Text>
        <Text style={styles.metricText}>Uptime: {formatUptime(process.uptime || 0)}</Text>
      </View>
      
      <View style={styles.processActions}>
        <TouchableOpacity
          style={[styles.actionButton, styles.startButton]}
          onPress={() => startProcess(process.name)}
          disabled={process.status === 'running'}
        >
          <Text style={styles.actionButtonText}>Start</Text>
        </TouchableOpacity>
        
        <TouchableOpacity
          style={[styles.actionButton, styles.stopButton]}
          onPress={() => stopProcess(process.name)}
          disabled={process.status === 'stopped'}
        >
          <Text style={styles.actionButtonText}>Stop</Text>
        </TouchableOpacity>
      </View>
    </View>
  );

  const formatUptime = (seconds) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  };

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="dark-content" backgroundColor="#f8f9fa" />
      
      <View style={styles.header}>
        <Text style={styles.headerTitle}>GProc Mobile</Text>
        <View style={[styles.connectionStatus, { backgroundColor: isConnected ? '#4CAF50' : '#F44336' }]}>
          <Text style={styles.connectionText}>{isConnected ? 'Connected' : 'Disconnected'}</Text>
        </View>
      </View>

      <ScrollView
        style={styles.scrollView}
        refreshControl={
          <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
        }
      >
        {processes.length === 0 ? (
          <View style={styles.emptyState}>
            <Text style={styles.emptyStateText}>No processes found</Text>
            <Text style={styles.emptyStateSubtext}>Pull down to refresh</Text>
          </View>
        ) : (
          processes.map((process, index) => (
            <ProcessCard key={index} process={process} />
          ))
        )}
      </ScrollView>
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f8f9fa',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: 16,
    backgroundColor: '#ffffff',
    borderBottomWidth: 1,
    borderBottomColor: '#e9ecef',
  },
  headerTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#212529',
  },
  connectionStatus: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 12,
  },
  connectionText: {
    color: '#ffffff',
    fontSize: 12,
    fontWeight: '600',
  },
  scrollView: {
    flex: 1,
    padding: 16,
  },
  processCard: {
    backgroundColor: '#ffffff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  processHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  processName: {
    fontSize: 18,
    fontWeight: '600',
    color: '#212529',
  },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 4,
    borderRadius: 8,
  },
  statusText: {
    color: '#ffffff',
    fontSize: 12,
    fontWeight: '600',
  },
  processCommand: {
    fontSize: 14,
    color: '#6c757d',
    marginBottom: 12,
  },
  processMetrics: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  metricText: {
    fontSize: 12,
    color: '#495057',
  },
  processActions: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  actionButton: {
    paddingHorizontal: 24,
    paddingVertical: 8,
    borderRadius: 8,
    minWidth: 80,
    alignItems: 'center',
  },
  startButton: {
    backgroundColor: '#28a745',
  },
  stopButton: {
    backgroundColor: '#dc3545',
  },
  actionButtonText: {
    color: '#ffffff',
    fontSize: 14,
    fontWeight: '600',
  },
  emptyState: {
    alignItems: 'center',
    justifyContent: 'center',
    paddingVertical: 64,
  },
  emptyStateText: {
    fontSize: 18,
    color: '#6c757d',
    marginBottom: 8,
  },
  emptyStateSubtext: {
    fontSize: 14,
    color: '#adb5bd',
  },
});

export default App;
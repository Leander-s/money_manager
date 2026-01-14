import { StatusBar } from 'expo-status-bar';
import { useState } from 'react';
import { SafeAreaView, StyleSheet } from 'react-native';

import LoginScreen from './src/screens/LoginScreen';
import MainScreen from './src/screens/MainScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import { theme } from './src/theme';

type Screen = 'login' | 'register' | 'main';

export default function App() {
  const [screen, setScreen] = useState<Screen>('login');
  const [displayName, setDisplayName] = useState('friend');

  const handleLogin = (name: string) => {
    setDisplayName(name);
    setScreen('main');
  };

  const handleRegister = (name: string) => {
    setDisplayName(name);
    setScreen('main');
  };

  const handleLogout = () => {
    setScreen('login');
  };

  return (
    <SafeAreaView style={styles.container}>
      {screen === 'login' ? (
        <LoginScreen onLogin={handleLogin} onGoToRegister={() => setScreen('register')} />
      ) : null}

      {screen === 'register' ? (
        <RegisterScreen onRegister={handleRegister} onGoToLogin={() => setScreen('login')} />
      ) : null}

      {screen === 'main' ? <MainScreen displayName={displayName} onLogout={handleLogout} /> : null}
      <StatusBar style="light" />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
});

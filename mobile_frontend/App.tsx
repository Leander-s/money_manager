import { StatusBar } from 'expo-status-bar';
import { useState } from 'react';
import { StyleSheet } from 'react-native';
import { SafeAreaProvider } from 'react-native-safe-area-context';

import LoginScreen from './src/screens/LoginScreen';
import MainScreen from './src/screens/MainScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import VerifyScreen from './src/screens/VerifyScreen';
import { fetchCurrentUser, login, register, AuthToken, User } from './src/api';
import { theme } from './src/theme';

type Screen = 'login' | 'register' | 'main' | 'verify';

export default function App() {
    const [screen, setScreen] = useState<Screen>('login');
    const [authToken, setAuthToken] = useState<AuthToken | null>(null);
    const [currentUser, setCurrentUser] = useState<User | null>(null);

    const handleLogin = async (email: string, password: string) => {
        try {
            const token = await login({ email, password });
            const user = await fetchCurrentUser(token);
            setAuthToken(token);
            setCurrentUser(user);
            setScreen('main');
            return { ok: true };
        } catch (error) {
            return { ok: false, error: getErrorMessage(error) };
        }
    };

    const handleRegister = async (username: string, email: string, password: string) => {
        try {
            await register({ username: username.trim() || null, email, password });
            setScreen('verify');
            return { ok: true };
        } catch (error) {
            return { ok: false, error: getErrorMessage(error) };
        }
    };

    const handleLogout = () => {
        setAuthToken(null);
        setCurrentUser(null);
        setScreen('login');
    };

    return (
        <SafeAreaProvider style={styles.container}>
            {screen === 'login' ? (
                <LoginScreen onLogin={handleLogin} onGoToRegister={() => setScreen('register')} />
            ) : null}

            {screen === 'register' ? (
                <RegisterScreen onRegister={handleRegister} onGoToLogin={() => setScreen('login')} />
            ) : null}

            {screen === 'verify' ? (
                <VerifyScreen onGoToLogin={() => setScreen('login')} />
            ) : null}

            {screen === 'main' && currentUser && authToken ? (
                <MainScreen user={currentUser} token={authToken} onLogout={handleLogout} />
            ) : null}
            <StatusBar style="light" />
        </SafeAreaProvider>
    );
}

function getErrorMessage(error: unknown) {
    if (error instanceof Error) {
        return error.message;
    }
    return 'Something went wrong';
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: theme.colors.background,
    },
});

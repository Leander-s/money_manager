import { StatusBar } from 'expo-status-bar';
import { useEffect, useState } from 'react';
import { Linking, StyleSheet } from 'react-native';
import { SafeAreaProvider } from 'react-native-safe-area-context';

import LoginScreen from './src/screens/LoginScreen';
import MainScreen from './src/screens/MainScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import VerifyScreen from './src/screens/VerifyScreen';
import ForgotPasswordScreen from './src/screens/ForgotPasswordScreen';
import ResetEmailSentScreen from './src/screens/ResetEmailSentScreen';
import ResetPasswordScreen from './src/screens/ResetPasswordScreen';
import {
    fetchCurrentUser,
    login,
    register,
    resetPassword,
    resetPasswordWithToken,
    AuthToken,
    User,
} from './src/api';
import { theme } from './src/theme';

type Screen =
    | 'login'
    | 'register'
    | 'main'
    | 'verify'
    | 'forgot'
    | 'reset-sent'
    | 'reset-password';

export default function App() {
    const [screen, setScreen] = useState<Screen>('login');
    const [authToken, setAuthToken] = useState<AuthToken | null>(null);
    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [resetEmail, setResetEmail] = useState<string | null>(null);
    const [resetToken, setResetToken] = useState<string | null>(null);

    useEffect(() => {
        const handleUrl = (url: string | null) => {
            const token = getResetTokenFromUrl(url);
            if (!token) {
                return;
            }
            setResetToken(token);
            setScreen('reset-password');
        };

        Linking.getInitialURL()
            .then(handleUrl)
            .catch(() => null);

        const subscription = Linking.addEventListener('url', ({ url }) => handleUrl(url));
        return () => {
            subscription.remove();
        };
    }, []);

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
        setResetEmail(null);
        setResetToken(null);
    };

    const handleRequestPasswordReset = async (email: string) => {
        try {
            await resetPassword({ email });
            setResetEmail(email);
            setScreen('reset-sent');
            return { ok: true };
        } catch (error) {
            return { ok: false, error: getErrorMessage(error) };
        }
    };

    const handleResetPassword = async (token: string, password: string) => {
        try {
            await resetPasswordWithToken({ token, password });
            return { ok: true };
        } catch (error) {
            return { ok: false, error: getErrorMessage(error) };
        }
    };

    const goToLogin = () => {
        setScreen('login');
        setResetEmail(null);
        setResetToken(null);
    };

    return (
        <SafeAreaProvider style={styles.container}>
            {screen === 'login' ? (
                <LoginScreen
                    onLogin={handleLogin}
                    onGoToRegister={() => setScreen('register')}
                    onGoToForgotPassword={() => setScreen('forgot')}
                />
            ) : null}

            {screen === 'register' ? (
                <RegisterScreen onRegister={handleRegister} onGoToLogin={goToLogin} />
            ) : null}

            {screen === 'verify' ? (
                <VerifyScreen onGoToLogin={goToLogin} />
            ) : null}

            {screen === 'forgot' ? (
                <ForgotPasswordScreen
                    onRequestReset={handleRequestPasswordReset}
                    onGoToLogin={goToLogin}
                />
            ) : null}

            {screen === 'reset-sent' ? (
                <ResetEmailSentScreen email={resetEmail ?? undefined} onGoToLogin={goToLogin} />
            ) : null}

            {screen === 'reset-password' ? (
                <ResetPasswordScreen
                    token={resetToken}
                    onResetPassword={handleResetPassword}
                    onGoToLogin={goToLogin}
                />
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

function getResetTokenFromUrl(url: string | null): string | null {
    if (!url) {
        return null;
    }

    try {
        const parsed = new URL(url);
        const segments = parsed.pathname.split('/').filter(Boolean);
        const resetIndex = segments.findIndex((segment) => segment === 'reset-password');
        if (resetIndex !== -1 && segments[resetIndex + 1]) {
            return decodeURIComponent(segments[resetIndex + 1]);
        }
        const token = parsed.searchParams.get('token');
        return token ? decodeURIComponent(token) : null;
    } catch (error) {
        return null;
    }
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: theme.colors.background,
    },
});

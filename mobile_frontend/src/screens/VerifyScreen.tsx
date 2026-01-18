import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import { theme } from '../theme';

type VerifyScreenProps = {
  onGoToLogin: () => void;
};

export default function VerifyScreen({ onGoToLogin }: VerifyScreenProps) {
    return (
        <View style={styles.container}>
            <View style={styles.card}>
                <Text style={styles.title}>Verify Your Email</Text>
                <Text style={styles.subtitle}>A verification link has been sent to your email address. Please check your inbox and click the link to verify your account.</Text>

                <Pressable style={styles.button} onPress={onGoToLogin}>
                    <Text style={styles.buttonText}>Back to Login</Text>
                </Pressable>
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: theme.colors.background,
    padding: 16,
  },
  card: {
    width: '100%',
    maxWidth: 400,
    backgroundColor: theme.colors.card,
    borderRadius: 8,
    padding: 24,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 5,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: theme.colors.textPrimary,
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 16,
    color: theme.colors.textMuted,
    marginBottom: 24,
  },
  button: {
    backgroundColor: theme.colors.accent,
    paddingVertical: 12,
    borderRadius: 6,
    alignItems: 'center',
  },
  buttonText: {
    color: theme.colors.textPrimary,
    fontSize: 16,
    fontWeight: '600',
  },
});

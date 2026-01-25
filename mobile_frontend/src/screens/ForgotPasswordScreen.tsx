import React, { useRef, useState } from 'react';
import { Pressable, StyleSheet, Text, TextInput, View } from 'react-native';

import { theme } from '../theme';

type ForgotPasswordScreenProps = {
  onRequestReset: (email: string) => Promise<{ ok: boolean; error?: string }>;
  onGoToLogin: () => void;
};

export default function ForgotPasswordScreen({
  onRequestReset,
  onGoToLogin,
}: ForgotPasswordScreenProps) {
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const emailRef = useRef<TextInput>(null);

  const handleReset = async () => {
    if (loading) {
      return;
    }
    const trimmedEmail = email.trim();
    if (!trimmedEmail) {
      setError('Email is required.');
      return;
    }
    setLoading(true);
    setError(null);
    const result = await onRequestReset(trimmedEmail);
    if (result.ok) {
      return;
    }
    setError(result.error ?? 'Unable to send reset email.');
    setLoading(false);
  };

  return (
    <View style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Forgot your password?</Text>
        <Text style={styles.subtitle}>Enter your email and we will send a reset link.</Text>

        <View style={styles.form}>
          <TextInput
            autoCapitalize="none"
            autoComplete="email"
            keyboardType="email-address"
            placeholder="Email"
            placeholderTextColor={theme.colors.textMuted}
            returnKeyType="done"
            style={styles.input}
            value={email}
            onChangeText={setEmail}
            onSubmitEditing={handleReset}
            enablesReturnKeyAutomatically
            ref={emailRef}
          />
        </View>

        <Pressable
          disabled={loading}
          style={({ pressed }) => [
            styles.primaryButton,
            pressed && styles.primaryButtonPressed,
            loading && styles.primaryButtonDisabled,
          ]}
          onPress={handleReset}
        >
          <Text style={styles.primaryButtonText}>
            {loading ? 'Sending...' : 'Send reset email'}
          </Text>
        </Pressable>

        {error ? <Text style={styles.errorText}>{error}</Text> : null}

        <Pressable style={styles.linkButton} onPress={onGoToLogin}>
          <Text style={styles.linkButtonText}>Back to login</Text>
        </Pressable>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    paddingHorizontal: 24,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: theme.colors.background,
  },
  card: {
    width: '100%',
    maxWidth: 360,
    padding: 28,
    borderRadius: theme.radii.md,
    backgroundColor: theme.colors.card,
    borderWidth: 1,
    borderColor: theme.colors.cardBorder,
    ...theme.shadow,
  },
  title: {
    fontSize: 24,
    fontWeight: '700',
    color: theme.colors.textPrimary,
    marginBottom: 6,
  },
  subtitle: {
    fontSize: 16,
    color: theme.colors.textMuted,
    marginBottom: 20,
  },
  form: {
    gap: 12,
    marginBottom: 20,
  },
  input: {
    borderWidth: 1,
    borderColor: theme.colors.inputBorder,
    borderRadius: theme.radii.sm,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 16,
    color: theme.colors.textPrimary,
    backgroundColor: theme.colors.inputBackground,
  },
  primaryButton: {
    backgroundColor: theme.colors.accent,
    paddingVertical: 14,
    borderRadius: theme.radii.sm,
    alignItems: 'center',
    marginBottom: 14,
  },
  primaryButtonPressed: {
    opacity: 0.9,
  },
  primaryButtonDisabled: {
    opacity: 0.7,
  },
  primaryButtonText: {
    color: theme.colors.accentText,
    fontSize: 16,
    fontWeight: '700',
  },
  errorText: {
    color: theme.colors.danger,
    marginBottom: 12,
    textAlign: 'center',
  },
  linkButton: {
    alignItems: 'center',
    paddingVertical: 4,
  },
  linkButtonText: {
    color: theme.colors.link,
    fontSize: 14,
    fontWeight: '600',
  },
});

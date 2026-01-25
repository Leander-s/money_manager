import React, { useRef, useState } from 'react';
import { Pressable, StyleSheet, Text, TextInput, View } from 'react-native';

import { theme } from '../theme';

type ResetPasswordScreenProps = {
  token: string | null;
  onResetPassword: (token: string, password: string) => Promise<{ ok: boolean; error?: string }>;
  onGoToLogin: () => void;
};

export default function ResetPasswordScreen({
  token,
  onResetPassword,
  onGoToLogin,
}: ResetPasswordScreenProps) {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const confirmRef = useRef<TextInput>(null);

  const handleReset = async () => {
    if (loading || success) {
      return;
    }
    if (!token) {
      setError('Reset link is invalid or has expired.');
      return;
    }
    if (!password || !confirmPassword) {
      setError('Please enter and confirm your new password.');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords do not match.');
      return;
    }

    setLoading(true);
    setError(null);
    const result = await onResetPassword(token, password);
    if (result.ok) {
      setSuccess(true);
      setLoading(false);
      return;
    }
    setError(result.error ?? 'Unable to reset password.');
    setLoading(false);
  };

  return (
    <View style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Reset your password</Text>
        <Text style={styles.subtitle}>Choose a new password for your account.</Text>

        <View style={styles.form}>
          <TextInput
            autoComplete="password"
            placeholder="New password"
            placeholderTextColor={theme.colors.textMuted}
            secureTextEntry
            returnKeyType="next"
            style={styles.input}
            value={password}
            onChangeText={(value) => {
              setPassword(value);
              if (success) {
                setSuccess(false);
              }
            }}
            onSubmitEditing={() => confirmRef.current?.focus()}
            enablesReturnKeyAutomatically
          />
          <TextInput
            autoComplete="password"
            placeholder="Confirm new password"
            placeholderTextColor={theme.colors.textMuted}
            secureTextEntry
            returnKeyType="done"
            style={styles.input}
            value={confirmPassword}
            onChangeText={(value) => {
              setConfirmPassword(value);
              if (success) {
                setSuccess(false);
              }
            }}
            onSubmitEditing={handleReset}
            enablesReturnKeyAutomatically
            ref={confirmRef}
          />
        </View>

        <Pressable
          disabled={loading || success}
          style={({ pressed }) => [
            styles.primaryButton,
            pressed && styles.primaryButtonPressed,
            (loading || success) && styles.primaryButtonDisabled,
          ]}
          onPress={handleReset}
        >
          <Text style={styles.primaryButtonText}>
            {loading ? 'Resetting...' : 'Reset password'}
          </Text>
        </Pressable>

        {error ? <Text style={styles.errorText}>{error}</Text> : null}
        {success ? (
          <Text style={styles.successText}>Password reset successful. You can log in now.</Text>
        ) : null}

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
  successText: {
    color: theme.colors.accent,
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

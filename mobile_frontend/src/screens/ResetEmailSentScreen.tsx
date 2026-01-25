import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import { theme } from '../theme';

type ResetEmailSentScreenProps = {
  email?: string;
  onGoToLogin: () => void;
};

export default function ResetEmailSentScreen({ email, onGoToLogin }: ResetEmailSentScreenProps) {
  const message = email
    ? `We sent a password reset link to ${email}. Check your inbox to continue.`
    : 'We sent a password reset email. Check your inbox to continue.';

  return (
    <View style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Email sent</Text>
        <Text style={styles.subtitle}>{message}</Text>

        <Pressable
          style={({ pressed }) => [styles.primaryButton, pressed && styles.primaryButtonPressed]}
          onPress={onGoToLogin}
        >
          <Text style={styles.primaryButtonText}>Back to login</Text>
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
  primaryButton: {
    backgroundColor: theme.colors.accent,
    paddingVertical: 14,
    borderRadius: theme.radii.sm,
    alignItems: 'center',
  },
  primaryButtonPressed: {
    opacity: 0.9,
  },
  primaryButtonText: {
    color: theme.colors.accentText,
    fontSize: 16,
    fontWeight: '700',
  },
});

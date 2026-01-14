import React from 'react';
import { Pressable, StyleSheet, Text, View } from 'react-native';

import { theme } from '../theme';

type MainScreenProps = {
  displayName: string;
  onLogout: () => void;
};

export default function MainScreen({ displayName, onLogout }: MainScreenProps) {
  return (
    <View style={styles.container}>
      <View style={styles.card}>
        <Text style={styles.title}>Main App</Text>
        <Text style={styles.subtitle}>Welcome, {displayName}.</Text>

        <View style={styles.infoPanel}>
          <Text style={styles.infoTitle}>You are in</Text>
          <Text style={styles.infoBody}>This is the only page for now.</Text>
        </View>

        <Pressable
          style={({ pressed }) => [styles.secondaryButton, pressed && styles.secondaryButtonPressed]}
          onPress={onLogout}
        >
          <Text style={styles.secondaryButtonText}>Log out</Text>
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
    color: theme.colors.accent,
    marginBottom: 20,
  },
  infoPanel: {
    borderRadius: theme.radii.sm,
    backgroundColor: theme.colors.inputBackground,
    padding: 16,
    borderWidth: 1,
    borderColor: theme.colors.inputBorder,
    marginBottom: 20,
  },
  infoTitle: {
    fontSize: 13,
    textTransform: 'uppercase',
    color: theme.colors.textMuted,
    marginBottom: 6,
    letterSpacing: 0.6,
  },
  infoBody: {
    fontSize: 16,
    color: theme.colors.textPrimary,
  },
  secondaryButton: {
    alignItems: 'center',
    paddingVertical: 12,
    borderRadius: theme.radii.sm,
    backgroundColor: theme.colors.danger,
  },
  secondaryButtonPressed: {
    opacity: 0.9,
  },
  secondaryButtonText: {
    color: theme.colors.dangerText,
    fontWeight: '700',
  },
});

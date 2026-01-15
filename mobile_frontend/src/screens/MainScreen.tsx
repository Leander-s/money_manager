import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, TextInput, View } from 'react-native';

import { AuthToken, createMoneyEntry, fetchLatestMoneyEntry, MoneyEntry, User } from '../api';
import { theme } from '../theme';

type MainScreenProps = {
  user: User;
  token: AuthToken;
  onLogout: () => void;
};

export default function MainScreen({ user, token, onLogout }: MainScreenProps) {
  const [balance, setBalance] = React.useState('');
  const [ratio, setRatio] = React.useState('');
  const [currentBudget, setCurrentBudget] = React.useState<number | null>(null);
  const [lastBalance, setLastBalance] = React.useState<number | null>(null);
  const [loading, setLoading] = React.useState(false);
  const [submitting, setSubmitting] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const ratioInputRef = React.useRef<TextInput>(null);

  const decimalRegex = /^\d*\.?\d*$/;

  const formatRatioPercent = React.useCallback((value: number) => {
    const percentValue = value * 100;
    if (!Number.isFinite(percentValue)) {
      return '';
    }
    if (Number.isInteger(percentValue)) {
      return String(percentValue);
    }
    return percentValue.toFixed(2).replace(/\.?0+$/, '');
  }, []);

  const applyLatestEntry = React.useCallback(
    (entry: MoneyEntry | null) => {
      if (!entry) {
        setCurrentBudget(null);
        setLastBalance(null);
        return;
      }
      setCurrentBudget(entry.budget);
      setLastBalance(entry.balance);
      setRatio(formatRatioPercent(entry.ratio));
    },
    [formatRatioPercent]
  );

  const handleBalanceChange = (value: string) => {
    if (decimalRegex.test(value)) {
      setBalance(value);
    }
  };

  const handleRatioChange = (value: string) => {
    if (!decimalRegex.test(value)) {
      return;
    }
    if (value === '' || value === '.') {
      setRatio(value);
      return;
    }
    const numericValue = Number(value);
    if (Number.isNaN(numericValue) || numericValue <= 100) {
      setRatio(value);
    }
  };

  const handleSubmit = React.useCallback(async () => {
    if (submitting) {
      return;
    }
    const parsedBalance = Number.parseFloat(balance);
    const parsedRatio = Number.parseFloat(ratio);
    if (!Number.isFinite(parsedBalance)) {
      setError('Please enter a valid balance.');
      return;
    }
    if (!Number.isFinite(parsedRatio) || parsedRatio < 0 || parsedRatio > 100) {
      setError('Ratio must be between 0 and 100.');
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      const entry = await createMoneyEntry(token, {
        balance: parsedBalance,
        ratio: parsedRatio / 100,
      });
      setBalance('');
      applyLatestEntry(entry);
    } catch (submitError) {
      setError(getErrorMessage(submitError));
    } finally {
      setSubmitting(false);
    }
  }, [applyLatestEntry, balance, ratio, submitting, token]);
  const handleBalanceSubmit = () => ratioInputRef.current?.focus();

  React.useEffect(() => {
    let cancelled = false;
    const loadLatestEntry = async () => {
      setLoading(true);
      setError(null);
      try {
        const entry = await fetchLatestMoneyEntry(token);
        if (!cancelled) {
          applyLatestEntry(entry);
        }
      } catch (loadError) {
        if (!cancelled) {
          setError(getErrorMessage(loadError));
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };
    loadLatestEntry();
    return () => {
      cancelled = true;
    };
  }, [token, applyLatestEntry]);

  const budgetLabel = loading
    ? 'Loading...'
    : currentBudget === null
      ? 'No data'
      : currentBudget.toFixed(2);
  const lastBalanceLabel = loading
    ? 'Loading...'
    : lastBalance === null
      ? 'No data'
      : lastBalance.toFixed(2);

  return (
    <ScrollView contentContainerStyle={styles.container} style={styles.scrollView}>
      <View style={styles.card}>
        <View style={styles.topRow}>
          <Text style={styles.emailText}>{user.email}</Text>
          <Pressable
            style={({ pressed }) => [styles.secondaryButton, pressed && styles.secondaryButtonPressed]}
            onPress={onLogout}
          >
            <Text style={styles.secondaryButtonText}>Log out</Text>
          </Pressable>
        </View>

        <Text style={styles.subtitle}>Logged in as {user.username || user.email}</Text>

        <View style={styles.inputRow}>
          <View style={styles.inputGroup}>
            <Text style={styles.label}>Balance</Text>
            <TextInput
              value={balance}
              onChangeText={handleBalanceChange}
              onSubmitEditing={handleBalanceSubmit}
              placeholder="0.00"
              placeholderTextColor={theme.colors.textMuted}
              keyboardType="decimal-pad"
              returnKeyType="next"
              style={styles.input}
            />
          </View>

          <View style={styles.inputGroup}>
            <Text style={styles.label}>% to Budget</Text>
            <TextInput
              ref={ratioInputRef}
              value={ratio}
              onChangeText={handleRatioChange}
              onSubmitEditing={handleSubmit}
              placeholder="0-100%"
              placeholderTextColor={theme.colors.textMuted}
              keyboardType="decimal-pad"
              returnKeyType="done"
              style={styles.input}
            />
          </View>

          <Pressable
            disabled={submitting}
            style={({ pressed }) => [
              styles.primaryButton,
              pressed && styles.primaryButtonPressed,
              submitting && styles.primaryButtonDisabled,
            ]}
            onPress={handleSubmit}
          >
            <Text style={styles.primaryButtonText}>
              {submitting ? 'Submitting...' : 'Submit'}
            </Text>
          </Pressable>
        </View>

        <Text style={styles.budgetText}>Current budget: {budgetLabel}</Text>
        <Text style={styles.lastBalanceText}>Last balance: {lastBalanceLabel}</Text>
        {error ? <Text style={styles.errorText}>{error}</Text> : null}
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  scrollView: {
    flex: 1,
    backgroundColor: theme.colors.background,
  },
  container: {
    flexGrow: 1,
    paddingHorizontal: 24,
    paddingVertical: 24,
    justifyContent: 'center',
    alignItems: 'center',
  },
  card: {
    width: '100%',
    maxWidth: 360,
    padding: 24,
    borderRadius: theme.radii.md,
    backgroundColor: theme.colors.card,
    borderWidth: 1,
    borderColor: theme.colors.cardBorder,
    ...theme.shadow,
  },
  topRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 12,
  },
  emailText: {
    fontSize: 12,
    color: theme.colors.textMuted,
  },
  subtitle: {
    fontSize: 20,
    fontWeight: '600',
    color: theme.colors.accent,
    marginBottom: 16,
  },
  inputRow: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    marginBottom: 8,
  },
  inputGroup: {
    flex: 1,
    marginRight: 12,
  },
  label: {
    fontSize: 13,
    color: theme.colors.textMuted,
    marginBottom: 8,
  },
  input: {
    borderRadius: theme.radii.sm,
    backgroundColor: theme.colors.inputBackground,
    borderWidth: 1,
    borderColor: theme.colors.inputBorder,
    paddingHorizontal: 12,
    paddingVertical: 10,
    color: theme.colors.textPrimary,
    fontSize: 15,
  },
  primaryButton: {
    alignItems: 'center',
    paddingVertical: 12,
    paddingHorizontal: 12,
    borderRadius: theme.radii.sm,
    backgroundColor: theme.colors.accent,
  },
  primaryButtonPressed: {
    opacity: 0.9,
  },
  primaryButtonText: {
    color: theme.colors.accentText,
    fontWeight: '700',
  },
  secondaryButton: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: theme.radii.sm,
    backgroundColor: theme.colors.danger,
  },
  secondaryButtonPressed: {
    opacity: 0.9,
  },
  secondaryButtonText: {
    color: theme.colors.dangerText,
    fontWeight: '600',
  },
  budgetText: {
    marginTop: 16,
    fontSize: 22,
    fontWeight: '700',
    color: theme.colors.textPrimary,
    textAlign: 'center',
  },
  lastBalanceText: {
    marginTop: 6,
    fontSize: 14,
    color: theme.colors.textMuted,
    textAlign: 'center',
  },
  primaryButtonDisabled: {
    opacity: 0.7,
  },
  errorText: {
    color: theme.colors.danger,
    marginTop: 12,
    textAlign: 'center',
  },
});

function getErrorMessage(error: unknown) {
  if (error instanceof Error) {
    return error.message;
  }
  return 'Something went wrong';
}

import React from 'react';
import { Pressable, ScrollView, StyleSheet, Text, TextInput, View } from 'react-native';

import {
    AuthToken,
    createMoneyEntry,
    deleteMoneyEntry,
    fetchMoneyEntries,
    MoneyEntry,
    updateMoneyEntry,
    User,
} from '../api';
import HistoryEntryCard from '../components/HistoryEntryCard';
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
    const [historyEntries, setHistoryEntries] = React.useState<MoneyEntry[]>([]);
    const [sidebarOpen, setSidebarOpen] = React.useState(false);
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

    const formatEntryDate = React.useCallback((value: string) => {
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) {
            return value;
        }
        return date.toLocaleString();
    }, []);

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
            setHistoryEntries((previous) => [
                entry,
                ...previous.filter((item) => item.id !== entry.id),
            ]);
        } catch (submitError) {
            setError(getErrorMessage(submitError));
        } finally {
            setSubmitting(false);
        }
    }, [applyLatestEntry, balance, ratio, submitting, token]);
    const handleBalanceSubmit = () => ratioInputRef.current?.focus();

    React.useEffect(() => {
        let cancelled = false;
        const loadEntries = async () => {
            setLoading(true);
            setError(null);
            try {
                let entries = await fetchMoneyEntries(token);
                if (cancelled) {
                    return;
                }
                if (entries === null) {
                    entries = []
                }
                setHistoryEntries(entries);
                applyLatestEntry(entries[0] ?? null);
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
        loadEntries();
        return () => {
            cancelled = true;
        };
    }, [token, applyLatestEntry]);

    const handleUpdateEntry = React.useCallback(
        async (id: number, updatedBalance: number, updatedRatio: number) => {
            setError(null);
            try {
                let updatedEntries = await updateMoneyEntry(token, {
                    id,
                    balance: updatedBalance,
                    ratio: updatedRatio,
                });
                if (updatedEntries === null) {
                    console.warn('Received null entries from API');
                }
                if (updatedEntries === null) {
                    updatedEntries = [];
                }
                console.log('Updated entries after update:', updatedEntries);
                setLastBalance(
                    updatedEntries.length > 0 ? updatedEntries[0].balance : null
                );
                setHistoryEntries(updatedEntries);
                setCurrentBudget(updatedEntries.length > 0 ? updatedEntries[0].budget : null);
                return true;
            } catch (updateError) {
                setError(getErrorMessage(updateError));
                return false;
            }
        },
        [applyLatestEntry, token]
    );

    const handleDeleteEntry = React.useCallback(
        async (id: number) => {
            setError(null);
            try {
                let moneyEntries = await deleteMoneyEntry(token, id);
                if (moneyEntries === null) {
                    console.warn('Received null entries from API');
                }
                if (moneyEntries === null) {
                    moneyEntries = [];
                }
                console.log('Updated entries after deletion:', moneyEntries);
                setLastBalance(
                    moneyEntries.length > 0 ? moneyEntries[0].balance : null
                );
                setHistoryEntries(moneyEntries);
                setCurrentBudget(moneyEntries.length > 0 ? moneyEntries[0].budget : null);
                return true;
            } catch (deleteError) {
                setError(getErrorMessage(deleteError));
                return false;
            }
        },
        [applyLatestEntry, token]
    );

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

    const displayName = user.username || user.email;
    const secondaryEmail = user.username ? user.email : null;

    return (
        <View style={styles.screen}>
            <ScrollView contentContainerStyle={styles.container} style={styles.scrollView}>
                <View style={styles.card}>
                    <View style={styles.topRow}>
                        <Text style={styles.title}>Money Manager</Text>
                        <Pressable
                            style={({ pressed }) => [styles.menuButton, pressed && styles.menuButtonPressed]}
                            onPress={() => setSidebarOpen(true)}
                        >
                            <Text style={styles.menuButtonText}>Account</Text>
                        </Pressable>
                    </View>

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

                    <View style={styles.summaryCard}>
                        <Text style={styles.summaryLabel}>Current budget</Text>
                        <Text style={styles.summaryValue}>{budgetLabel}</Text>
                        <Text style={styles.summaryLabel}>Last balance</Text>
                        <Text style={styles.summarySubvalue}>{lastBalanceLabel}</Text>
                    </View>
                    <View style={styles.historyContainer}>
                        <Text style={styles.historyTitle}>History</Text>
                        {loading ? (
                            <Text style={styles.historyPlaceholder}>Loading...</Text>
                        ) : historyEntries.length === 0 ? (
                            <Text style={styles.historyPlaceholder}>No history yet.</Text>
                        ) : (
                            historyEntries.map((entry) => (
                                <HistoryEntryCard
                                    key={entry.id}
                                    entry={entry}
                                    formatEntryDate={formatEntryDate}
                                    formatRatioPercent={formatRatioPercent}
                                    onUpdate={handleUpdateEntry}
                                    onDelete={handleDeleteEntry}
                                />
                            )
                            )
                        )}
                    </View>
                    {error ? <Text style={styles.errorText}>{error}</Text> : null}
                </View>
            </ScrollView>
            {sidebarOpen ? (
                <View style={styles.sidebarOverlay}>
                    <Pressable
                        style={styles.sidebarBackdrop}
                        onPress={() => setSidebarOpen(false)}
                    />
                    <View style={styles.sidebar}>
                        <Text style={styles.sidebarName}>{displayName}</Text>
                        {secondaryEmail ? (
                            <Text style={styles.sidebarEmail}>{secondaryEmail}</Text>
                        ) : null}
                        <Pressable
                            style={({ pressed }) => [
                                styles.secondaryButton,
                                pressed && styles.secondaryButtonPressed,
                            ]}
                            onPress={() => {
                                setSidebarOpen(false);
                                onLogout();
                            }}
                        >
                            <Text style={styles.secondaryButtonText}>Log out</Text>
                        </Pressable>
                    </View>
                </View>
            ) : null}
        </View>
    );
}

const styles = StyleSheet.create({
    screen: {
        flex: 1,
        backgroundColor: theme.colors.background,
    },
    scrollView: {
        flex: 1,
        backgroundColor: theme.colors.background,
    },
    container: {
        flexGrow: 1,
        paddingHorizontal: 24,
        paddingVertical: 24,
        alignItems: 'center',
    },
    card: {
        width: '100%',
        maxWidth: 520,
        alignSelf: 'center',
        padding: 24,
    },
    topRow: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: 16,
    },
    title: {
        fontSize: 20,
        fontWeight: '700',
        color: theme.colors.textPrimary,
    },
    menuButton: {
        paddingHorizontal: 12,
        paddingVertical: 8,
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.inputBackground,
        borderWidth: 1,
        borderColor: theme.colors.inputBorder,
    },
    menuButtonPressed: {
        opacity: 0.8,
    },
    menuButtonText: {
        color: theme.colors.textPrimary,
        fontWeight: '600',
        fontSize: 13,
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
    summaryCard: {
        marginTop: 24,
        paddingVertical: 18,
        paddingHorizontal: 16,
        borderRadius: theme.radii.md,
        backgroundColor: 'rgba(255, 255, 255, 0.04)',
        borderWidth: 1,
        borderColor: 'rgba(255, 255, 255, 0.12)',
    },
    summaryLabel: {
        fontSize: 13,
        textTransform: 'uppercase',
        letterSpacing: 1,
        color: theme.colors.textMuted,
        textAlign: 'center',
    },
    summaryValue: {
        marginTop: 6,
        marginBottom: 14,
        fontSize: 34,
        fontWeight: '800',
        color: theme.colors.textPrimary,
        textAlign: 'center',
    },
    summarySubvalue: {
        marginTop: 6,
        fontSize: 20,
        fontWeight: '700',
        color: theme.colors.textPrimary,
        textAlign: 'center',
    },
    historyContainer: {
        marginTop: 16,
        width: '100%',
    },
    historyTitle: {
        fontSize: 16,
        fontWeight: '700',
        color: theme.colors.textPrimary,
        marginBottom: 8,
    },
    historyPlaceholder: {
        color: theme.colors.textMuted,
        textAlign: 'center',
    },
    historyItem: {
        padding: 12,
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.inputBackground,
        borderWidth: 1,
        borderColor: theme.colors.inputBorder,
        marginBottom: 10,
    },
    historyDate: {
        color: theme.colors.textMuted,
        fontSize: 12,
        marginBottom: 6,
    },
    historyRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        marginTop: 2,
    },
    historyLabel: {
        color: theme.colors.textMuted,
        fontSize: 12,
    },
    historyValue: {
        color: theme.colors.textPrimary,
        fontSize: 13,
        fontWeight: '600',
    },
    primaryButtonDisabled: {
        opacity: 0.7,
    },
    errorText: {
        color: theme.colors.danger,
        marginTop: 12,
        textAlign: 'center',
    },
    sidebarOverlay: {
        position: 'absolute',
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        flexDirection: 'row',
    },
    sidebarBackdrop: {
        flex: 1,
        backgroundColor: 'rgba(0, 0, 0, 0.45)',
    },
    sidebar: {
        width: 260,
        backgroundColor: theme.colors.background,
        paddingTop: 48,
        paddingHorizontal: 20,
        borderLeftWidth: 1,
        borderLeftColor: theme.colors.cardBorder,
    },
    sidebarName: {
        fontSize: 18,
        fontWeight: '700',
        color: theme.colors.textPrimary,
    },
    sidebarEmail: {
        marginTop: 6,
        fontSize: 12,
        color: theme.colors.textMuted,
        marginBottom: 16,
    },
});

function getErrorMessage(error: unknown) {
    if (error instanceof Error) {
        return error.message;
    }
    return 'Something went wrong';
}

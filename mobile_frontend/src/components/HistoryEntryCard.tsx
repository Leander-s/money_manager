import React from 'react';
import { Pressable, StyleSheet, Text, TextInput, View } from 'react-native';

import { MoneyEntry } from '../api';
import { theme } from '../theme';

type HistoryEntryCardProps = {
    entry: MoneyEntry;
    formatEntryDate: (value: string) => string;
    formatRatioPercent: (value: number) => string;
    onUpdate: (id: number, balance: number, ratio: number) => Promise<boolean>;
    onDelete: (id: number) => Promise<boolean>;
};

export default function HistoryEntryCard({
    entry,
    formatEntryDate,
    formatRatioPercent,
    onUpdate,
    onDelete,
}: HistoryEntryCardProps) {
    const [isEditing, setIsEditing] = React.useState(false);
    const [draftBalance, setDraftBalance] = React.useState(entry.balance.toFixed(2));
    const [draftRatio, setDraftRatio] = React.useState(formatRatioPercent(entry.ratio));
    const [validationError, setValidationError] = React.useState<string | null>(null);
    const [saving, setSaving] = React.useState(false);
    const [deleting, setDeleting] = React.useState(false);
    const ratioInputRef = React.useRef<TextInput>(null);

    const decimalRegex = /^\d*\.?\d*$/;

    React.useEffect(() => {
        if (isEditing) {
            return;
        }
        setDraftBalance(entry.balance.toFixed(2));
        setDraftRatio(formatRatioPercent(entry.ratio));
        setValidationError(null);
    }, [entry.balance, entry.ratio, formatRatioPercent, isEditing]);

    const handleEditPress = () => {
        if (saving || deleting) {
            return;
        }
        setIsEditing(true);
        setValidationError(null);
    };

    const handleCancel = () => {
        if (saving || deleting) {
            return;
        }
        setIsEditing(false);
        setDraftBalance(entry.balance.toFixed(2));
        setDraftRatio(formatRatioPercent(entry.ratio));
        setValidationError(null);
    };

    const handleBalanceChange = (value: string) => {
        if (decimalRegex.test(value)) {
            setDraftBalance(value);
        }
    };

    const handleRatioChange = (value: string) => {
        if (!decimalRegex.test(value)) {
            return;
        }
        if (value === '' || value === '.') {
            setDraftRatio(value);
            return;
        }
        const numericValue = Number(value);
        if (Number.isNaN(numericValue) || numericValue <= 100) {
            setDraftRatio(value);
        }
    };

    const handleSave = async () => {
        if (saving || deleting) {
            return;
        }
        const parsedBalance = Number.parseFloat(draftBalance);
        const parsedRatio = Number.parseFloat(draftRatio);
        if (!Number.isFinite(parsedBalance)) {
            setValidationError('Please enter a valid balance.');
            return;
        }
        if (!Number.isFinite(parsedRatio) || parsedRatio < 0 || parsedRatio > 100) {
            setValidationError('Ratio must be between 0 and 100.');
            return;
        }
        setSaving(true);
        setValidationError(null);
        const success = await onUpdate(entry.id, parsedBalance, parsedRatio / 100);
        setSaving(false);
        if (success) {
            setIsEditing(false);
        }
    };

    const handleDelete = async () => {
        if (saving || deleting) {
            return;
        }
        setDeleting(true);
        const success = await onDelete(entry.id);
        setDeleting(false);
        if (success) {
            setIsEditing(false);
        }
    };

    return (
        <View style={styles.card}>
            <View style={styles.headerRow}>
                <Text style={styles.dateText}>{formatEntryDate(entry.created_at)}</Text>
                <View style={styles.actionRow}>
                    <Pressable
                        disabled={saving || deleting || isEditing}
                        style={({ pressed }) => [
                            styles.actionButton,
                            pressed && styles.actionButtonPressed,
                            (saving || deleting || isEditing) && styles.actionButtonDisabled,
                        ]}
                        onPress={handleEditPress}
                    >
                        <Text style={styles.editButtonText}>Edit</Text>
                    </Pressable>
                    <Pressable
                        disabled={saving || deleting}
                        style={({ pressed }) => [
                            styles.deleteButton,
                            pressed && styles.actionButtonPressed,
                            (saving || deleting) && styles.actionButtonDisabled,
                        ]}
                        onPress={handleDelete}
                    >
                        <Text style={styles.deleteButtonText}>
                            {deleting ? 'Deleting...' : 'Delete'}
                        </Text>
                    </Pressable>
                </View>
            </View>

            {isEditing ? (
                <>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Balance</Text>
                        <TextInput
                            value={draftBalance}
                            onChangeText={handleBalanceChange}
                            onSubmitEditing={() => ratioInputRef.current?.focus()}
                            placeholder="0.00"
                            placeholderTextColor={theme.colors.textMuted}
                            keyboardType="decimal-pad"
                            returnKeyType="next"
                            style={styles.historyInput}
                        />
                    </View>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Budget</Text>
                        <Text style={styles.historyValue}>{entry.budget.toFixed(2)}</Text>
                    </View>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Ratio</Text>
                        <View style={styles.ratioInputWrap}>
                            <TextInput
                                ref={ratioInputRef}
                                value={draftRatio}
                                onChangeText={handleRatioChange}
                                onSubmitEditing={handleSave}
                                placeholder="0-100"
                                placeholderTextColor={theme.colors.textMuted}
                                keyboardType="decimal-pad"
                                returnKeyType="done"
                                style={styles.historyInput}
                            />
                            <Text style={styles.ratioSuffix}>%</Text>
                        </View>
                    </View>
                    {validationError ? (
                        <Text style={styles.validationError}>{validationError}</Text>
                    ) : null}
                    <View style={styles.editActionsRow}>
                        <Pressable
                            disabled={saving || deleting}
                            style={({ pressed }) => [
                                styles.saveButton,
                                pressed && styles.actionButtonPressed,
                                (saving || deleting) && styles.actionButtonDisabled,
                            ]}
                            onPress={handleSave}
                        >
                            <Text style={styles.saveButtonText}>
                                {saving ? 'Saving...' : 'Save'}
                            </Text>
                        </Pressable>
                        <Pressable
                            disabled={saving || deleting}
                            style={({ pressed }) => [
                                styles.cancelButton,
                                pressed && styles.actionButtonPressed,
                                (saving || deleting) && styles.actionButtonDisabled,
                            ]}
                            onPress={handleCancel}
                        >
                            <Text style={styles.cancelButtonText}>Cancel</Text>
                        </Pressable>
                    </View>
                </>
            ) : (
                <>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Balance</Text>
                        <Text style={styles.historyValue}>{entry.balance.toFixed(2)}</Text>
                    </View>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Budget</Text>
                        <Text style={styles.historyValue}>{entry.budget.toFixed(2)}</Text>
                    </View>
                    <View style={styles.historyRow}>
                        <Text style={styles.historyLabel}>Ratio</Text>
                        <Text style={styles.historyValue}>
                            {formatRatioPercent(entry.ratio)}%
                        </Text>
                    </View>
                </>
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    card: {
        padding: 12,
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.inputBackground,
        borderWidth: 1,
        borderColor: theme.colors.inputBorder,
        marginBottom: 10,
    },
    headerRow: {
        flexDirection: 'row',
        alignItems: 'center',
        justifyContent: 'space-between',
        marginBottom: 6,
    },
    dateText: {
        color: theme.colors.textMuted,
        fontSize: 12,
        flex: 1,
        marginRight: 8,
    },
    actionRow: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    actionButton: {
        paddingHorizontal: 10,
        paddingVertical: 6,
        borderRadius: theme.radii.sm,
        backgroundColor: 'rgba(147, 197, 253, 0.15)',
        borderWidth: 1,
        borderColor: 'rgba(147, 197, 253, 0.35)',
        marginLeft: 8,
    },
    deleteButton: {
        paddingHorizontal: 10,
        paddingVertical: 6,
        borderRadius: theme.radii.sm,
        backgroundColor: 'rgba(239, 68, 68, 0.12)',
        borderWidth: 1,
        borderColor: 'rgba(239, 68, 68, 0.4)',
        marginLeft: 8,
    },
    actionButtonPressed: {
        opacity: 0.8,
    },
    actionButtonDisabled: {
        opacity: 0.6,
    },
    editButtonText: {
        color: theme.colors.link,
        fontWeight: '600',
        fontSize: 12,
    },
    deleteButtonText: {
        color: theme.colors.danger,
        fontWeight: '600',
        fontSize: 12,
    },
    historyRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
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
    historyInput: {
        minWidth: 90,
        textAlign: 'right',
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.inputBackground,
        borderWidth: 1,
        borderColor: theme.colors.inputBorder,
        paddingHorizontal: 10,
        paddingVertical: 6,
        color: theme.colors.textPrimary,
        fontSize: 13,
    },
    ratioInputWrap: {
        flexDirection: 'row',
        alignItems: 'center',
    },
    ratioSuffix: {
        color: theme.colors.textMuted,
        fontSize: 12,
        fontWeight: '600',
        marginLeft: 6,
    },
    validationError: {
        color: theme.colors.danger,
        marginTop: 6,
        fontSize: 12,
    },
    editActionsRow: {
        flexDirection: 'row',
        justifyContent: 'flex-end',
        marginTop: 10,
    },
    saveButton: {
        paddingHorizontal: 12,
        paddingVertical: 8,
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.accent,
        marginRight: 8,
    },
    saveButtonText: {
        color: theme.colors.accentText,
        fontWeight: '700',
        fontSize: 12,
    },
    cancelButton: {
        paddingHorizontal: 12,
        paddingVertical: 8,
        borderRadius: theme.radii.sm,
        backgroundColor: theme.colors.card,
        borderWidth: 1,
        borderColor: theme.colors.cardBorder,
    },
    cancelButtonText: {
        color: theme.colors.textPrimary,
        fontWeight: '600',
        fontSize: 12,
    },
});

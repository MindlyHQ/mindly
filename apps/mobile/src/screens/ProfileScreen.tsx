import React from 'react';
import { View, Text, StyleSheet, TouchableOpacity, ScrollView, Alert } from 'react-native';
import { useAuth } from '../../src/contexts/AuthContext';

export default function ProfileScreen({ navigation }: any) {
    const { user, logout } = useAuth();

    const handleLogout = async () => {
        Alert.alert(
            'Выход',
            'Вы уверены, что хотите выйти?',
            [
                { text: 'Отмена', style: 'cancel' },
                {
                    text: 'Выйти',
                    style: 'destructive',
                    onPress: async () => {
                        await logout();
                        navigation.navigate('Feed');
                    }
                }
            ]
        );
    };

    if (!user) {
        return (
            <View style={styles.container}>
                <Text>Загрузка профиля...</Text>
            </View>
        );
    }

    return (
        <ScrollView style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.title}>👤 Профиль</Text>
                <Text style={styles.username}>@{user.username}</Text>
            </View>

            <View style={styles.statsContainer}>
                <View style={styles.statCard}>
                    <Text style={styles.statValue}>{user.score}</Text>
                    <Text style={styles.statLabel}>Очки</Text>
                </View>
                <View style={styles.statCard}>
                    <Text style={styles.statValue}>{user.current_streak}</Text>
                    <Text style={styles.statLabel}>Дней подряд</Text>
                </View>
                <View style={styles.statCard}>
                    <Text style={styles.statValue}>{user.best_streak}</Text>
                    <Text style={styles.statLabel}>Рекорд</Text>
                </View>
            </View>

            <View style={styles.infoSection}>
                <View style={styles.infoCard}>
                    <Text style={styles.infoLabel}>Полное имя</Text>
                    <Text style={styles.infoValue}>{user.full_name || 'Не указано'}</Text>
                </View>

                <View style={styles.infoCard}>
                    <Text style={styles.infoLabel}>Email</Text>
                    <Text style={styles.infoValue}>{user.email}</Text>
                </View>

                <View style={styles.infoCard}>
                    <Text style={styles.infoLabel}>Дата регистрации</Text>
                    <Text style={styles.infoValue}>
                        {new Date(user.created_at).toLocaleDateString('ru-RU')}
                    </Text>
                </View>
            </View>

            <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
                <Text style={styles.logoutText}>🚪 Выйти из аккаунта</Text>
            </TouchableOpacity>
        </ScrollView>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f8f9fa',
    },
    header: {
        backgroundColor: '#4a6fa5',
        paddingVertical: 30,
        paddingHorizontal: 20,
        borderBottomLeftRadius: 20,
        borderBottomRightRadius: 20,
        marginBottom: 20,
    },
    title: {
        fontSize: 32,
        fontWeight: 'bold',
        color: 'white',
        marginBottom: 5,
    },
    username: {
        fontSize: 16,
        color: '#e2e8f0',
    },
    statsContainer: {
        flexDirection: 'row',
        justifyContent: 'space-around',
        paddingHorizontal: 15,
        marginBottom: 25,
    },
    statCard: {
        backgroundColor: 'white',
        padding: 15,
        borderRadius: 15,
        alignItems: 'center',
        minWidth: 90,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 4,
        elevation: 3,
    },
    statValue: {
        fontSize: 24,
        fontWeight: 'bold',
        color: '#4a6fa5',
    },
    statLabel: {
        fontSize: 12,
        color: '#718096',
        marginTop: 5,
    },
    infoSection: {
        paddingHorizontal: 15,
        marginBottom: 30,
    },
    infoCard: {
        backgroundColor: 'white',
        padding: 18,
        borderRadius: 12,
        marginBottom: 12,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 1 },
        shadowOpacity: 0.08,
        shadowRadius: 3,
        elevation: 2,
    },
    infoLabel: {
        fontSize: 14,
        color: '#718096',
        marginBottom: 6,
    },
    infoValue: {
        fontSize: 16,
        fontWeight: '600',
        color: '#2d3748',
    },
    logoutButton: {
        backgroundColor: '#e53e3e',
        marginHorizontal: 20,
        paddingVertical: 16,
        borderRadius: 12,
        alignItems: 'center',
        marginBottom: 30,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 3 },
        shadowOpacity: 0.15,
        shadowRadius: 5,
        elevation: 4,
    },
    logoutText: {
        color: 'white',
        fontSize: 16,
        fontWeight: '600',
    },
});
import React, { useState } from 'react';
import {
    View,
    Text,
    TextInput,
    Button,
    Alert,
    StyleSheet,
    ScrollView,
    ActivityIndicator,
    TouchableOpacity,
} from 'react-native';
import axios from 'axios';

// ВАЖНО: Для Android эмулятора используйте 'http://10.0.2.2:8080'
// Для iOS симулятора или реального устройства: 'http://localhost:8080' или ваш локальный IP
const API_BASE_URL = 'http://192.168.0.101:8081';

interface RegisterRequest {
    email: string;
    username: string;
    password: string;
    full_name?: string;
}

interface User {
    id: string;
    email: string;
    username: string;
    full_name?: string;
    score: number;
    current_streak: number;
    best_streak: number;
    created_at: string;
    updated_at: string;
}

export default function RegisterScreen() {
    const [formData, setFormData] = useState<RegisterRequest>({
        email: '',
        username: '',
        password: '',
        full_name: '',
    });
    const [loading, setLoading] = useState(false);
    const [registeredUser, setRegisteredUser] = useState<User | null>(null);
    const [serverStatus, setServerStatus] = useState<'checking' | 'online' | 'offline'>('checking');

    // Проверка связи с сервером при загрузке экрана
    React.useEffect(() => {
        checkServerHealth();
    }, []);

    const checkServerHealth = async () => {
        try {
            const response = await axios.get(`${API_BASE_URL}/health`, { timeout: 5000 });
            if (response.data.status === 'ok') {
                setServerStatus('online');
                console.log('✅ Сервер доступен:', response.data);
            } else {
                setServerStatus('offline');
            }
        } catch (error) {
            console.error('❌ Сервер недоступен:', error);
            setServerStatus('offline');
        }
    };

    const handleInputChange = (field: keyof RegisterRequest, value: string) => {
        setFormData(prev => ({ ...prev, [field]: value }));
    };

    const validateForm = (): boolean => {
        if (!formData.email.includes('@')) {
            Alert.alert('Ошибка', 'Введите корректный email');
            return false;
        }
        if (formData.username.length < 3) {
            Alert.alert('Ошибка', 'Имя пользователя должно быть не менее 3 символов');
            return false;
        }
        if (formData.password.length < 6) {
            Alert.alert('Ошибка', 'Пароль должен быть не менее 6 символов');
            return false;
        }
        return true;
    };

    const handleRegister = async () => {
        if (!validateForm()) return;

        setLoading(true);
        console.log('📤 Отправка запроса регистрации:', formData);

        try {
            const response = await axios.post(
                `${API_BASE_URL}/api/auth/register`,
                {
                    email: formData.email.trim(),
                    username: formData.username.trim(),
                    password: formData.password,
                    full_name: formData.full_name?.trim() || null,
                },
                {
                    timeout: 10000,
                    headers: { 'Content-Type': 'application/json' },
                }
            );

            console.log('✅ Ответ сервера:', response.data);

            if (response.data.status === 'success') {
                const user: User = response.data.data;
                setRegisteredUser(user);

                Alert.alert(
                    '🎉 Успешная регистрация!',
                    `Пользователь ${user.username} создан!\nID: ${user.id.substring(0, 8)}...`,
                    [{ text: 'OK' }]
                );

                // Очистка формы
                setFormData({
                    email: '',
                    username: '',
                    password: '',
                    full_name: '',
                });
            } else {
                throw new Error(response.data.error || 'Неизвестная ошибка сервера');
            }
        } catch (error: any) {
            console.error('❌ Ошибка регистрации:', error);

            let errorMessage = 'Неизвестная ошибка';
            if (error.response?.data?.error) {
                errorMessage = error.response.data.error;
            } else if (error.code === 'ECONNABORTED') {
                errorMessage = 'Таймаут запроса. Сервер не отвечает';
            } else if (error.message.includes('Network Error')) {
                errorMessage = 'Не удалось подключиться к серверу. Проверьте, что сервер запущен';
            }

            Alert.alert('❌ Ошибка регистрации', errorMessage);
        } finally {
            setLoading(false);
        }
    };

    const handleTestUser = () => {
        setFormData({
            email: 'test@mindly.com',
            username: 'mindly_user',
            password: 'password123',
            full_name: 'Test User',
        });
    };

    return (
        <ScrollView contentContainerStyle={styles.container}>
            {/* Заголовок */}
            <View style={styles.header}>
                <Text style={styles.title}>🧠 Mindly</Text>
                <Text style={styles.subtitle}>Регистрация пользователя</Text>

                {/* Индикатор статуса сервера */}
                <View style={[
                    styles.statusBadge,
                    serverStatus === 'online' ? styles.statusOnline :
                        serverStatus === 'offline' ? styles.statusOffline : styles.statusChecking
                ]}>
                    <Text style={styles.statusText}>
                        {serverStatus === 'online' ? '✅ Сервер онлайн' :
                            serverStatus === 'offline' ? '❌ Сервер офлайн' : '🔄 Проверка...'}
                    </Text>
                </View>
            </View>

            {/* Карточка формы */}
            <View style={styles.card}>
                <Text style={styles.sectionTitle}>📝 Форма регистрации</Text>

                <TextInput
                    style={styles.input}
                    placeholder="Email *"
                    value={formData.email}
                    onChangeText={(value) => handleInputChange('email', value)}
                    autoCapitalize="none"
                    keyboardType="email-address"
                    editable={!loading}
                />

                <TextInput
                    style={styles.input}
                    placeholder="Имя пользователя *"
                    value={formData.username}
                    onChangeText={(value) => handleInputChange('username', value)}
                    autoCapitalize="none"
                    editable={!loading}
                />

                <TextInput
                    style={styles.input}
                    placeholder="Пароль * (мин. 6 символов)"
                    value={formData.password}
                    onChangeText={(value) => handleInputChange('password', value)}
                    secureTextEntry
                    editable={!loading}
                />

                <TextInput
                    style={styles.input}
                    placeholder="Полное имя (необязательно)"
                    value={formData.full_name}
                    onChangeText={(value) => handleInputChange('full_name', value)}
                    editable={!loading}
                />

                <View style={styles.buttonContainer}>
                    {loading ? (
                        <ActivityIndicator size="large" color="#4a6fa5" />
                    ) : (
                        <>
                            <TouchableOpacity style={styles.primaryButton} onPress={handleRegister}>
                                <Text style={styles.buttonText}>Зарегистрироваться</Text>
                            </TouchableOpacity>

                            <TouchableOpacity style={styles.secondaryButton} onPress={handleTestUser}>
                                <Text style={styles.secondaryButtonText}>Заполнить тестовыми данными</Text>
                            </TouchableOpacity>

                            <TouchableOpacity style={styles.secondaryButton} onPress={checkServerHealth}>
                                <Text style={styles.secondaryButtonText}>Проверить соединение с сервером</Text>
                            </TouchableOpacity>
                        </>
                    )}
                </View>
            </View>

            {/* Карточка зарегистрированного пользователя */}
            {registeredUser && (
                <View style={[styles.card, styles.successCard]}>
                    <Text style={styles.sectionTitle}>✅ Зарегистрированный пользователь</Text>
                    <View style={styles.userInfo}>
                        <Text style={styles.userInfoLabel}>ID:</Text>
                        <Text style={styles.userInfoValue}>{registeredUser.id.substring(0, 16)}...</Text>

                        <Text style={styles.userInfoLabel}>Email:</Text>
                        <Text style={styles.userInfoValue}>{registeredUser.email}</Text>

                        <Text style={styles.userInfoLabel}>Username:</Text>
                        <Text style={styles.userInfoValue}>{registeredUser.username}</Text>

                        <Text style={styles.userInfoLabel}>Score:</Text>
                        <Text style={styles.userInfoValue}>{registeredUser.score}</Text>

                        <Text style={styles.userInfoLabel}>Registered:</Text>
                        <Text style={styles.userInfoValue}>
                            {new Date(registeredUser.created_at).toLocaleDateString('ru-RU')}
                        </Text>
                    </View>
                </View>
            )}

            {/* Информационная карточка */}
            <View style={styles.infoCard}>
                <Text style={styles.infoTitle}>ℹ️ Информация о Дне 2</Text>
                <Text style={styles.infoText}>
                    • Реализован полный цикл регистрации{'\n'}
                    • Go API + PostgreSQL + CORS{'\n'}
                    • React Native с TypeScript{'\n'}
                    • Валидация формы на клиенте{'\n'}
                    • Обработка ошибок сети{'\n'}
                    • Статус соединения в реальном времени
                </Text>
            </View>
        </ScrollView>
    );
}

const styles = StyleSheet.create({
    container: {
        flexGrow: 1,
        padding: 20,
        backgroundColor: '#f5f7fa',
    },
    header: {
        alignItems: 'center',
        marginBottom: 30,
    },
    title: {
        fontSize: 36,
        fontWeight: 'bold',
        color: '#2c5282',
        marginBottom: 5,
    },
    subtitle: {
        fontSize: 16,
        color: '#4a5568',
        marginBottom: 15,
    },
    statusBadge: {
        paddingHorizontal: 15,
        paddingVertical: 6,
        borderRadius: 20,
    },
    statusOnline: {
        backgroundColor: '#c6f6d5',
    },
    statusOffline: {
        backgroundColor: '#fed7d7',
    },
    statusChecking: {
        backgroundColor: '#feebc8',
    },
    statusText: {
        fontSize: 14,
        fontWeight: '600',
    },
    card: {
        backgroundColor: 'white',
        borderRadius: 12,
        padding: 20,
        marginBottom: 20,
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 6,
        elevation: 3,
    },
    successCard: {
        borderLeftWidth: 4,
        borderLeftColor: '#48bb78',
    },
    sectionTitle: {
        fontSize: 18,
        fontWeight: '600',
        color: '#2d3748',
        marginBottom: 20,
    },
    input: {
        borderWidth: 1,
        borderColor: '#e2e8f0',
        borderRadius: 8,
        padding: 12,
        marginBottom: 15,
        fontSize: 16,
        backgroundColor: '#f8fafc',
    },
    buttonContainer: {
        marginTop: 10,
    },
    primaryButton: {
        backgroundColor: '#4a6fa5',
        borderRadius: 8,
        padding: 15,
        alignItems: 'center',
        marginBottom: 10,
    },
    buttonText: {
        color: 'white',
        fontSize: 16,
        fontWeight: '600',
    },
    secondaryButton: {
        backgroundColor: '#edf2f7',
        borderRadius: 8,
        padding: 12,
        alignItems: 'center',
        marginBottom: 10,
    },
    secondaryButtonText: {
        color: '#4a5568',
        fontSize: 14,
    },
    userInfo: {
        marginTop: 10,
    },
    userInfoLabel: {
        fontWeight: '600',
        color: '#4a5568',
        marginTop: 8,
    },
    userInfoValue: {
        color: '#2d3748',
        marginBottom: 4,
    },
    infoCard: {
        backgroundColor: '#ebf8ff',
        borderRadius: 8,
        padding: 15,
        borderLeftWidth: 4,
        borderLeftColor: '#4299e1',
    },
    infoTitle: {
        fontSize: 16,
        fontWeight: '600',
        color: '#2b6cb0',
        marginBottom: 10,
    },
    infoText: {
        color: '#4a5568',
        fontSize: 14,
        lineHeight: 20,
    },
});
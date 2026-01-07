import React, { useEffect, useState } from 'react';
import {
    View,
    Text,
    FlatList,
    Image,
    StyleSheet,
    Dimensions,
    ActivityIndicator,
    TouchableOpacity,
    Platform
} from 'react-native';
import { useAuth } from '../../src/contexts/AuthContext';

const { width } = Dimensions.get('window');

// ⚠️ ВАЖНО: ПОДСТАВЬ СВОЙ IP КОМПЬЮТЕРА СЮДА!
const COMPUTER_IP = '192.168.0.160'; // ← ИЗМЕНИ ЭТУ СТРОКУ! Используй IP из ipconfig

// УНИВЕРСАЛЬНОЕ РЕШЕНИЕ ДЛЯ РЕАЛЬНОГО ТЕЛЕФОНА
const getApiUrl = () => {
    if (__DEV__) {
        // Для реального телефона всегда используем IP компьютера
        return `http://${COMPUTER_IP}:8081/api/feed`;
    }
    return 'https://api.твой-домен.com/api/feed'; // Для продакшена
};

const API_URL = getApiUrl();

type VideoItem = {
    id: string;
    title: string;
    description: string;
    video_url: string;
    thumbnail_url: string;
    duration_sec: number;
    tags: string[];
    author: {
        full_name: string;
        expertise_area: string;
        trust_tier: string;
    };
};

export default function FeedScreen({ navigation }: any) {
    const [videos, setVideos] = useState<VideoItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [retryCount, setRetryCount] = useState(0);
    const { user } = useAuth(); // Получаем данные пользователя

    useEffect(() => {
        loadVideos();
    }, [retryCount]);

    const loadVideos = async () => {
        try {
            setLoading(true);
            setError(null);

            console.log('📱 Платформа:', Platform.OS);
            console.log('🔗 Загружаем видео с:', API_URL);
            console.log('📡 Точный URL:', `${API_URL}?limit=10`);

            // Добавляем таймаут для запроса (5 секунд)
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), 5000);

            const response = await fetch(`${API_URL}?limit=10`, {
                signal: controller.signal,
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                }
            });

            clearTimeout(timeoutId);

            console.log('✅ Статус ответа:', response.status);
            console.log('📊 Заголовки ответа:', Object.fromEntries(response.headers.entries()));

            const data = await response.json();
            console.log('🎬 Получено видео:', data.data?.length || 0);

            if (data.success) {
                setVideos(data.data || []);
            } else {
                setError(`Ошибка API: ${data.error || 'Неизвестная ошибка'}`);
            }
        } catch (err: any) {
            console.error('❌ Ошибка загрузки:', err.message);
            console.error('🔧 Тип ошибки:', err.name);

            let errorMessage = 'Не удалось подключиться к серверу';

            if (err.name === 'AbortError') {
                errorMessage = 'Таймаут подключения (сервер не ответил за 5 секунд)';
            } else if (err.message.includes('Network request failed')) {
                errorMessage = 'Сетевая ошибка. Проверьте:';
            }

            setError(`${errorMessage}\nURL: ${API_URL}`);
        } finally {
            setLoading(false);
        }
    };

    // Изменено: теперь ведет на профиль ИЛИ регистрацию
    const navigateToProfileOrRegister = () => {
        if (user) {
            navigation.navigate('Profile');
        } else {
            navigation.navigate('Register');
        }
    };

    const retryWithDifferentUrl = () => {
        // Показываем варианты для разных платформ
        alert(
            `Попробуйте изменить URL в коде:\n\n` +
            `• iOS симулятор: http://localhost:8080/api/feed\n` +
            `• Android эмулятор: http://10.0.2.2:8080/api/feed\n` +
            `• Реальный телефон: http://ВАШ_IP:8080/api/feed\n\n` +
            `Текущий URL: ${API_URL}`
        );
        setRetryCount(prev => prev + 1);
    };

    const testApiInBrowser = () => {
        // Показываем инструкцию для проверки в браузере
        alert(
            `Откройте браузер на компьютере и перейдите по адресу:\n\n` +
            `http://localhost:8080/api/feed\n\n` +
            `Если видите JSON с видео — API работает.\n` +
            `Если нет — проверьте Go-сервер.`
        );
    };

    // Функция для лайков (только для авторизованных)
    const handleLikePress = (videoId: string) => {
        if (!user) {
            alert('Войдите в аккаунт, чтобы ставить лайки');
            return;
        }
        // TODO: Добавить логику лайков с user.id
        console.log(`Пользователь ${user.id} лайкнул видео ${videoId}`);
    };

    if (loading) {
        return (
            <View style={styles.center}>
                <ActivityIndicator size="large" color="#4a6fa5" />
                <Text style={styles.loadingText}>Загружаем образовательные видео...</Text>
                <Text style={styles.urlHint}>Подключаемся к: {API_URL}</Text>
            </View>
        );
    }

    if (error) {
        return (
            <View style={styles.center}>
                <Text style={styles.errorTitle}>⚠️ Ошибка подключения</Text>
                <Text style={styles.error}>{error}</Text>

                <View style={styles.hintBox}>
                    <Text style={styles.hintTitle}>Возможные причины:</Text>
                    <Text style={styles.hint}>1. Go-сервер не запущен</Text>
                    <Text style={styles.hint}>2. Неправильный URL для платформы</Text>
                    <Text style={styles.hint}>3. Проблемы с сетью или CORS</Text>
                </View>

                <View style={styles.buttonContainer}>
                    <TouchableOpacity style={styles.button} onPress={loadVideos}>
                        <Text style={styles.buttonText}>🔄 Повторить попытку</Text>
                    </TouchableOpacity>

                    <TouchableOpacity style={[styles.button, styles.secondaryButton]} onPress={retryWithDifferentUrl}>
                        <Text style={styles.buttonText}>🔧 Изменить URL</Text>
                    </TouchableOpacity>

                    <TouchableOpacity style={[styles.button, styles.infoButton]} onPress={testApiInBrowser}>
                        <Text style={styles.buttonText}>🌐 Проверить API в браузере</Text>
                    </TouchableOpacity>

                    {/* Изменено: теперь показывает "Профиль" или "Регистрация" */}
                    <TouchableOpacity style={[styles.button, styles.registerButton]} onPress={navigateToProfileOrRegister}>
                        <Text style={styles.buttonText}>{user ? '👤 Профиль' : '👤 Регистрация'}</Text>
                    </TouchableOpacity>
                </View>

                <Text style={styles.platformInfo}>
                    Платформа: {Platform.OS} | URL: {API_URL}
                </Text>
            </View>
        );
    }

    const renderVideo = ({ item }: { item: VideoItem }) => (
        <View style={styles.videoContainer}>
            <Image
                source={{ uri: item.thumbnail_url || 'https://via.placeholder.com/300x500' }}
                style={styles.thumbnail}
                resizeMode="cover"
            />
            <View style={styles.videoInfo}>
                <View style={styles.headerRow}>
                    <Text style={styles.title}>{item.title}</Text>
                    <Text style={styles.duration}>{item.duration_sec} сек</Text>
                </View>

                <View style={styles.authorRow}>
                    <Text style={styles.authorName}>{item.author.full_name}</Text>
                    <View style={[styles.badge,
                    item.author.trust_tier === 'gold' ? styles.goldBadge :
                        item.author.trust_tier === 'silver' ? styles.silverBadge :
                            styles.bronzeBadge
                    ]}>
                        <Text style={styles.badgeText}>{item.author.trust_tier}</Text>
                    </View>
                </View>

                <Text style={styles.expertise}>{item.author.expertise_area}</Text>

                <View style={styles.tagsContainer}>
                    {item.tags.slice(0, 3).map((tag, index) => (
                        <View key={index} style={styles.tag}>
                            <Text style={styles.tagText}>#{tag}</Text>
                        </View>
                    ))}
                </View>

                <Text style={styles.description} numberOfLines={2}>{item.description}</Text>

                {/* Кнопка лайка (показывает статус авторизации) */}
                <TouchableOpacity
                    style={styles.likeButton}
                    onPress={() => handleLikePress(item.id)}
                >
                    <Text style={styles.likeText}>
                        {user ? '❤️ Лайк' : '🔒 Войдите, чтобы лайкнуть'}
                    </Text>
                </TouchableOpacity>
            </View>
        </View>
    );

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <Text style={styles.headerTitle}>🎓 LearnStream</Text>
                <Text style={styles.headerSubtitle}>Образовательные видео в формате TikTok</Text>
                <Text style={styles.urlDisplay}>API: {API_URL.replace('http://', '')}</Text>
                {/* Показываем статус пользователя в заголовке */}
                {user && (
                    <Text style={styles.userStatus}>
                        👤 Вы вошли как {user.username}
                    </Text>
                )}
            </View>

            <FlatList
                data={videos}
                renderItem={renderVideo}
                keyExtractor={(item) => item.id}
                showsVerticalScrollIndicator={false}
                contentContainerStyle={styles.listContent}
                refreshing={loading}
                onRefresh={loadVideos}
                ListEmptyComponent={
                    <View style={styles.emptyContainer}>
                        <Text style={styles.emptyText}>Нет доступных видео</Text>
                        <TouchableOpacity onPress={loadVideos}>
                            <Text style={styles.emptyLink}>Обновить</Text>
                        </TouchableOpacity>
                    </View>
                }
            />

            {/* Изменено: теперь показывает разную иконку в зависимости от авторизации */}
            <TouchableOpacity style={styles.fab} onPress={navigateToProfileOrRegister}>
                <Text style={styles.fabText}>{user ? '👤' : '📝'}</Text>
            </TouchableOpacity>

            <View style={styles.debugInfo}>
                <Text style={styles.debugText}>📱 {Platform.OS.toUpperCase()}</Text>
                <Text style={styles.debugText}>🎬 {videos.length} видео</Text>
                {/* Добавляем информацию о пользователе в debug */}
                <Text style={styles.debugText}>
                    {user ? `👤 ${user.username.substring(0, 8)}...` : '👤 Не авторизован'}
                </Text>
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#f8f9fa',
    },
    center: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: 20,
        backgroundColor: '#f8f9fa',
    },
    header: {
        backgroundColor: '#4a6fa5',
        paddingTop: 50,
        paddingBottom: 15,
        paddingHorizontal: 15,
        borderBottomLeftRadius: 20,
        borderBottomRightRadius: 20,
        marginBottom: 10,
    },
    headerTitle: {
        fontSize: 28,
        fontWeight: 'bold',
        color: 'white',
        marginBottom: 5,
    },
    headerSubtitle: {
        fontSize: 14,
        color: '#e2e8f0',
        marginBottom: 5,
    },
    urlDisplay: {
        fontSize: 10,
        color: '#a0c4ff',
        fontFamily: 'monospace',
    },
    userStatus: {
        fontSize: 12,
        color: '#cbd5e0',
        marginTop: 5,
        fontStyle: 'italic',
    },
    loadingText: {
        marginTop: 15,
        fontSize: 16,
        color: '#666',
    },
    urlHint: {
        fontSize: 12,
        color: '#888',
        marginTop: 10,
        fontFamily: 'monospace',
    },
    errorTitle: {
        fontSize: 20,
        fontWeight: 'bold',
        color: '#e53e3e',
        marginBottom: 10,
    },
    error: {
        color: '#e53e3e',
        fontSize: 14,
        marginBottom: 20,
        textAlign: 'center',
        fontFamily: 'monospace',
        backgroundColor: '#fed7d7',
        padding: 10,
        borderRadius: 8,
        borderWidth: 1,
        borderColor: '#fc8181',
    },
    hintBox: {
        backgroundColor: '#ebf8ff',
        padding: 15,
        borderRadius: 10,
        marginBottom: 20,
        borderWidth: 1,
        borderColor: '#bee3f8',
    },
    hintTitle: {
        fontWeight: 'bold',
        color: '#2c5282',
        marginBottom: 8,
    },
    hint: {
        color: '#4a5568',
        fontSize: 14,
        marginBottom: 5,
    },
    buttonContainer: {
        width: '100%',
        alignItems: 'center',
    },
    button: {
        backgroundColor: '#4a6fa5',
        paddingVertical: 12,
        paddingHorizontal: 25,
        borderRadius: 25,
        marginTop: 10,
        minWidth: 250,
        alignItems: 'center',
        flexDirection: 'row',
        justifyContent: 'center',
    },
    secondaryButton: {
        backgroundColor: '#718096',
    },
    infoButton: {
        backgroundColor: '#38a169',
    },
    registerButton: {
        backgroundColor: '#d69e2e',
    },
    buttonText: {
        color: 'white',
        fontWeight: '600',
        fontSize: 16,
        marginLeft: 8,
    },
    platformInfo: {
        marginTop: 20,
        fontSize: 12,
        color: '#718096',
        fontFamily: 'monospace',
    },
    listContent: {
        paddingBottom: 20,
    },
    videoContainer: {
        backgroundColor: 'white',
        marginHorizontal: 15,
        marginBottom: 15,
        borderRadius: 15,
        overflow: 'hidden',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 2 },
        shadowOpacity: 0.1,
        shadowRadius: 6,
        elevation: 3,
    },
    thumbnail: {
        width: '100%',
        height: (width - 30) * 1.5,
        backgroundColor: '#e2e8f0',
    },
    videoInfo: {
        padding: 15,
    },
    headerRow: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: 8,
    },
    title: {
        fontSize: 18,
        fontWeight: '700',
        color: '#2d3748',
        flex: 1,
        marginRight: 10,
    },
    duration: {
        fontSize: 13,
        color: '#718096',
        backgroundColor: '#f7fafc',
        paddingHorizontal: 8,
        paddingVertical: 3,
        borderRadius: 12,
    },
    authorRow: {
        flexDirection: 'row',
        alignItems: 'center',
        marginBottom: 5,
    },
    authorName: {
        fontSize: 15,
        fontWeight: '600',
        color: '#4a5568',
        marginRight: 10,
    },
    badge: {
        paddingHorizontal: 8,
        paddingVertical: 2,
        borderRadius: 10,
    },
    goldBadge: {
        backgroundColor: '#f6e05e',
    },
    silverBadge: {
        backgroundColor: '#cbd5e0',
    },
    bronzeBadge: {
        backgroundColor: '#ed8936',
    },
    badgeText: {
        fontSize: 11,
        fontWeight: '700',
        color: '#2d3748',
        textTransform: 'uppercase',
    },
    expertise: {
        fontSize: 14,
        color: '#4a6fa5',
        marginBottom: 10,
        fontStyle: 'italic',
    },
    tagsContainer: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        marginBottom: 12,
    },
    tag: {
        backgroundColor: '#ebf8ff',
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
        marginRight: 8,
        marginBottom: 5,
    },
    tagText: {
        fontSize: 12,
        color: '#2b6cb0',
        fontWeight: '500',
    },
    description: {
        fontSize: 14,
        color: '#4a5568',
        lineHeight: 20,
        marginBottom: 10,
    },
    likeButton: {
        backgroundColor: '#f7fafc',
        paddingVertical: 8,
        paddingHorizontal: 15,
        borderRadius: 20,
        alignSelf: 'flex-start',
        borderWidth: 1,
        borderColor: '#e2e8f0',
    },
    likeText: {
        fontSize: 14,
        fontWeight: '500',
        color: '#4a5568',
    },
    fab: {
        position: 'absolute',
        bottom: 20,
        right: 20,
        backgroundColor: '#4a6fa5',
        width: 56,
        height: 56,
        borderRadius: 28,
        justifyContent: 'center',
        alignItems: 'center',
        shadowColor: '#000',
        shadowOffset: { width: 0, height: 4 },
        shadowOpacity: 0.3,
        shadowRadius: 4,
        elevation: 6,
    },
    fabText: {
        fontSize: 24,
        color: 'white',
    },
    emptyContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        padding: 50,
    },
    emptyText: {
        fontSize: 18,
        color: '#718096',
        marginBottom: 10,
    },
    emptyLink: {
        fontSize: 16,
        color: '#4a6fa5',
        textDecorationLine: 'underline',
    },
    debugInfo: {
        position: 'absolute',
        top: 10,
        right: 10,
        backgroundColor: 'rgba(0,0,0,0.7)',
        paddingHorizontal: 10,
        paddingVertical: 5,
        borderRadius: 10,
        flexDirection: 'row',
        gap: 10,
    },
    debugText: {
        color: 'white',
        fontSize: 10,
        fontFamily: 'monospace',
    },
});
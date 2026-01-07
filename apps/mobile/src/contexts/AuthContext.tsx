import React, { createContext, useState, useEffect, useContext } from 'react';
import AsyncStorage from '@react-native-async-storage/async-storage';

export interface User {
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

export interface AuthContextType {
    user: User | null;
    isLoading: boolean;
    login: (userData: User) => Promise<void>;
    logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);
const USER_STORAGE_KEY = '@mindly_user';

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        loadUser();
    }, []);

    const loadUser = async () => {
        try {
            const userJson = await AsyncStorage.getItem(USER_STORAGE_KEY);
            if (userJson) {
                setUser(JSON.parse(userJson));
            }
        } catch (error) {
            console.error('Auth load error:', error);
        } finally {
            setIsLoading(false);
        }
    };

    const login = async (userData: User) => {
        try {
            await AsyncStorage.setItem(USER_STORAGE_KEY, JSON.stringify(userData));
            setUser(userData);
        } catch (error) {
            console.error('Auth login error:', error);
            throw error;
        }
    };

    const logout = async () => {
        try {
            await AsyncStorage.removeItem(USER_STORAGE_KEY);
            setUser(null);
        } catch (error) {
            console.error('Auth logout error:', error);
            throw error;
        }
    };

    return (
        <AuthContext.Provider value={{ user, isLoading, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within AuthProvider');
    }
    return context;
}
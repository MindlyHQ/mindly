import React from 'react';
import { AppRegistry } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';

// Импортируем экраны
import FeedScreen from './src/screens/FeedScreen';
import RegisterScreen from './src/screens/RegisterScreen';
import ProfileScreen from './src/screens/ProfileScreen';

// Импортируем AuthProvider
import { AuthProvider } from './src/contexts/AuthContext';

// Создаём навигатор
const Stack = createNativeStackNavigator();

// Главный компонент приложения
function App() {
  return (
    <AuthProvider>
      <NavigationContainer>
        <Stack.Navigator initialRouteName="Feed">
          {/* Экран ленты видео (главный) */}
          <Stack.Screen
            name="Feed"
            component={FeedScreen}
            options={{
              title: 'Mindly',
              headerStyle: {
                backgroundColor: '#4a6fa5',
              },
              headerTintColor: 'white',
              headerTitleStyle: {
                fontWeight: 'bold',
                fontSize: 20,
              },
            }}
          />

          {/* Экран регистрации */}
          <Stack.Screen
            name="Register"
            component={RegisterScreen}
            options={{
              title: 'Регистрация',
              headerStyle: {
                backgroundColor: '#4a6fa5',
              },
              headerTintColor: 'white',
            }}
          />

          {/* Экран профиля */}
          <Stack.Screen
            name="Profile"
            component={ProfileScreen}
            options={{
              title: 'Профиль',
              headerStyle: {
                backgroundColor: '#4a6fa5',
              },
              headerTintColor: 'white',
            }}
          />
        </Stack.Navigator>
      </NavigationContainer>
    </AuthProvider>
  );
}

// ✅ КРИТИЧЕСКИ ВАЖНО: Регистрируем компонент
AppRegistry.registerComponent('main', () => App);

// Экспортируем компонент
export default App;
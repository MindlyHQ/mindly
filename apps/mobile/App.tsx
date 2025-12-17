// App.tsx - С ЯВНОЙ РЕГИСТРАЦИЕЙ КОМПОНЕНТА
import React from 'react';
import { AppRegistry, View, Text } from 'react-native';

// 1. Определяем компонент
function App() {
  return (
    <View style={{
      flex: 1,
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: '#4a6fa5'
    }}>
      <Text style={{
        fontSize: 32,
        fontWeight: 'bold',
        color: 'white',
        marginBottom: 20
      }}>
        🧠 Mindly
      </Text>
      <Text style={{
        fontSize: 18,
        color: '#e2e8f0'
      }}>
        День 2: Регистрация работает!
      </Text>
    </View>
  );
}

// 2. ✅ КРИТИЧЕСКИ ВАЖНО: Регистрируем компонент
AppRegistry.registerComponent('main', () => App);

// 3. Экспортируем компонент
export default App;
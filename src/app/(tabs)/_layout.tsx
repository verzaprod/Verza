import { Tabs, } from "expo-router";
import { Icon, } from "@/components/ui/Icon";


export default function TabsLayout() {
  return (  

    <Tabs>  
      <Tabs.Screen  
        name="Home"  
        options={{
          tabBarLabel: 'Home',
          tabBarIcon: ({ color, size }) => (
            <Icon name="home" size={size} color={color} />
          ),
          headerShown: false,
        }}
      />

    </Tabs>
  )
}

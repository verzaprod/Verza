import { Tabs, Stack } from "expo-router";
import { Icon, } from "@/components/ui/Icon";


export default function TabsLayout() {
  return (  
    <>
      {/* <Stack /> */}
       <Tabs screenOptions={{ headerShown: false }}>  
         <Tabs.Screen  
           name="home"  
           options={{
             tabBarLabel: 'Home',
             tabBarIcon: ({ color, size }) => (
               <Icon name="home" size={size} color={color} />
             ),
             headerShown: false,
           }}
         />

       </Tabs>
    </>
  )
}

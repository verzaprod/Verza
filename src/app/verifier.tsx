import React, { useState } from "react";
import { SafeAreaView, View, TouchableOpacity } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useRouter } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import VerificationJobList from "@/components/verifiers/VerificationJobList";
import DocumentReviewModal from "@/components/verifiers/DocumentReviewModal";
import { useKYCStore } from "@/store/kycStore";
import { useTheme } from "@/theme/ThemeProvider";

const mockJobs = [
  { 
    id: "1", 
    requester: "James", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id1.jpg",
    avatar: "https://i.pravatar.cc/150?img=1"
  },
  { 
    id: "2", 
    requester: "Aliyat", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id2.jpg",
    avatar: "https://i.pravatar.cc/150?img=2"
  },
  { 
    id: "3", 
    requester: "Richard", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id3.jpg",
    avatar: "https://i.pravatar.cc/150?img=3"
  },
  { 
    id: "4", 
    requester: "Tunde", 
    doc: "Identity Card", 
    status: "pending",
    documentImage: "https://example.com/id4.jpg",
    avatar: "https://i.pravatar.cc/150?img=4"
  },
];

export default function VerifierDashboard() {
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const theme = useTheme();
  const setVerificationStatus = useKYCStore((state) => state.setVerificationStatus);
  
  const [jobs, setJobs] = useState(mockJobs);
  const [selectedJob, setSelectedJob] = useState(null);

  const handleApprove = (job) => {
    setVerificationStatus("verified");
    setJobs((prev) => prev.filter((j) => j.id !== job.id));
    setSelectedJob(null);
    
    if (jobs.length <= 1) {
      router.replace("/(tabs)/home");
    }
  };

  const handleReject = (job) => {
    setVerificationStatus("rejected");
    setJobs((prev) => prev.filter((j) => j.id !== job.id));
    setSelectedJob(null);
    
    if (jobs.length <= 1) {
      router.replace("/(tabs)/home");
    }
  };

  return (
    <SafeAreaView
      style={{
        flex: 1,
        backgroundColor: theme.colors.background,
        paddingTop: insets.top,
      }}
    >
      {/* Header */}
      <View style={{
        flexDirection: "row",
        justifyContent: "space-between",
        alignItems: "center",
        paddingHorizontal: 20,
        marginBottom: theme.spacing.lg,
      }}>
        <TouchableOpacity>
          <View style={{
            width: 64,
            height: 64,
            borderRadius: 32,
            backgroundColor: theme.colors.primaryGreen,
            alignItems: "center",
            justifyContent: "center",
            overflow: "hidden",
          }}>
            <Ionicons name="person" size={32} color="#FFFFFF" />
          </View>
        </TouchableOpacity>
        
        <View style={{ flexDirection: "row", gap: 16 }}>
          <TouchableOpacity>
            <Ionicons 
              name="notifications-outline" 
              size={28} 
              color={theme.colors.textSecondary} 
            />
          </TouchableOpacity>
          <TouchableOpacity>
            <Ionicons 
              name="search-outline" 
              size={28} 
              color={theme.colors.textSecondary} 
            />
          </TouchableOpacity>
        </View>
      </View>

      <VerificationJobList
        jobs={jobs}
        onJobPress={setSelectedJob}
      />
      
      <DocumentReviewModal
        visible={!!selectedJob}
        job={selectedJob}
        onClose={() => setSelectedJob(null)}
        onApprove={handleApprove}
        onReject={handleReject}
      />

      {/* Bottom Navigation */}
      <View 
        style={{
          position: "absolute",
          bottom: 0,
          left: 0,
          right: 0,
          height: 96,
          flexDirection: "row",
          justifyContent: "space-between",
          alignItems: "center",
          paddingHorizontal: 48,
          backgroundColor: theme.colors.background,
          paddingBottom: insets.bottom,
        }}
      >
        <TouchableOpacity style={{ width: 48, height: 48, alignItems: "center", justifyContent: "center" }}>
          <Ionicons 
            name="menu-outline" 
            size={32} 
            color={theme.colors.textSecondary} 
          />
        </TouchableOpacity>

        <TouchableOpacity 
          style={{
            width: 80,
            height: 80,
            borderRadius: 40,
            backgroundColor: theme.colors.primaryGreen,
            alignItems: "center",
            justifyContent: "center",
            marginTop: -40,
            shadowColor: theme.colors.primaryGreen,
            shadowOffset: { width: 0, height: 4 },
            shadowOpacity: 0.3,
            shadowRadius: 12,
            elevation: 8,
          }}
        >
          <Ionicons name="home" size={32} color="#FFFFFF" />
        </TouchableOpacity>

        <TouchableOpacity style={{ width: 48, height: 48, alignItems: "center", justifyContent: "center" }}>
          <Ionicons 
            name="person-outline" 
            size={32} 
            color={theme.colors.textSecondary} 
          />
        </TouchableOpacity>
      </View>
    </SafeAreaView> 
  );
} 
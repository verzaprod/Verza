import React, { useState } from "react";
import { SafeAreaView, View, Text, TouchableOpacity, FlatList, Modal } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useTheme } from "@/theme/ThemeProvider";
import { useKYCStore } from "@/store/kycStore";
import { useRouter } from "expo-router";

const mockJobs = [
  { id: "1", requester: "Alice", doc: "Passport", status: "available" },
  { id: "2", requester: "Bob", doc: "Driver License", status: "active" },
  { id: "3", requester: "Carol", doc: "ID Card", status: "completed" },
];

function JobDetailModal({ visible, job, onClose, onApprove, onReject }) {
  const theme = useTheme();
  return (
    <Modal visible={visible} transparent animationType="slide">
      <View style={{ flex: 1, justifyContent: "center", padding: 20, backgroundColor: "rgba(0,0,0,0.6)" }}>
        <View style={{ backgroundColor: theme.colors.background, borderRadius: 16, padding: 20 }}>
          <Text style={{ fontSize: 20, fontWeight: "bold", color: theme.colors.textPrimary, marginBottom: 12 }}>
            Job Detail
          </Text>
          <Text style={{ color: theme.colors.textSecondary, marginBottom: 8 }}>
            Requester: {job?.requester}
          </Text>
          <Text style={{ color: theme.colors.textSecondary, marginBottom: 20 }}>
            Document: {job?.doc}
          </Text>
          <TouchableOpacity
            style={{ backgroundColor: theme.colors.primaryGreen, padding: 14, borderRadius: 12, marginBottom: 12 }}
            onPress={() => { onApprove(job); onClose(); }}
          >
            <Text style={{ color: theme.colors.textPrimary, textAlign: "center", fontWeight: "bold" }}>Approve</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={{ backgroundColor: theme.colors.error, padding: 14, borderRadius: 12 }}
            onPress={() => { onReject(job); onClose(); }}
          >
            <Text style={{ color: theme.colors.secondaryAccent, textAlign: "center", fontWeight: "bold" }}>Reject</Text>
          </TouchableOpacity>
        </View>
      </View>
    </Modal>
  );
}

export default function VerifierDashboard() {
  const insets = useSafeAreaInsets();
  const theme = useTheme();
  const setVerificationStatus  = useKYCStore((state) => state.setVerificationStatus);
  const [jobs, setJobs] = useState(mockJobs);
  const [selectedJob, setSelectedJob] = useState(null);

  const router = useRouter();

  const handleApprove = (job) => {
    setVerificationStatus("verified");
    setJobs((prev) => prev.map((j) => (j.id === job.id ? { ...j, status: "completed" } : j)));
    router.replace("/(tabs)/home");
  };

  const handleReject = (job) => {
    setVerificationStatus("rejected");
    setJobs((prev) => prev.map((j) => (j.id === job.id ? { ...j, status: "completed" } : j)));
    router.replace("/(tabs)/home");
  };

  const renderJob = ({ item }) => (
    <TouchableOpacity
      style={{
        padding: 16,
        backgroundColor: theme.colors.background,
        borderRadius: 12,
        marginBottom: 12,
      }}
      onPress={() => setSelectedJob(item)}
    >
      <Text style={{ color: theme.colors.textPrimary, fontWeight: "600" }}>
        {item.requester}
      </Text>
      <Text style={{ color: theme.colors.textSecondary, fontSize: 14 }}>
        {item.doc} — {item.status}
      </Text>
    </TouchableOpacity>
  );

  return (
    <SafeAreaView
      style={{
        flex: 1,
        paddingTop: insets.top,
        backgroundColor: theme.colors.background,
        paddingHorizontal: 20,
      }}
    >
      <Text
        style={{
          fontSize: 24,
          fontWeight: "bold",
          color: theme.colors.textPrimary,
          marginVertical: 20,
          textAlign: "center",
        }}
      >
        Verifier Dashboard
      </Text>
      <FlatList data={jobs} renderItem={renderJob} keyExtractor={(item) => item.id} />
      <JobDetailModal
        visible={!!selectedJob}
        job={selectedJob}
        onClose={() => setSelectedJob(null)}
        onApprove={handleApprove}
        onReject={handleReject}
      />
    </SafeAreaView>
  );
}

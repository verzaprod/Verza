import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';

interface Rating {
  verifierId: string;
  rating: number;
  comment?: string;
  timestamp: string;
}

interface RatingStore {
  ratings: Rating[];
  currentRating: number;
  currentComment: string;
  isSubmitting: boolean;
  setCurrentRating: (rating: number) => void;
  setCurrentComment: (comment: string) => void;
  submitRating: (verifierId: string) => Promise<void>;
  loadRatings: () => Promise<void>;
  getRatingForVerifier: (verifierId: string) => Rating | null;
}

export const useRatingStore = create<RatingStore>((set, get) => ({
  ratings: [],
  currentRating: 0,
  currentComment: '',
  isSubmitting: false,

  setCurrentRating: (rating: number) => set({ currentRating: rating }),
  
  setCurrentComment: (comment: string) => set({ currentComment: comment }),

  submitRating: async (verifierId: string) => {
    const { currentRating, currentComment, ratings } = get();
    
    if (currentRating === 0) return;

    set({ isSubmitting: true });

    try {
      const newRating: Rating = {
        verifierId,
        rating: currentRating,
        comment: currentComment,
        timestamp: new Date().toISOString(),
      };

      const updatedRatings = [
        ...ratings.filter(r => r.verifierId !== verifierId),
        newRating
      ];

      await SecureStore.setItem('verifier_ratings', JSON.stringify(updatedRatings));
      
      set({ 
        ratings: updatedRatings,
        currentRating: 0,
        currentComment: '',
        isSubmitting: false
      });
    } catch (error) {
      console.error('Failed to save rating:', error);
      set({ isSubmitting: false });
    }
  },

  loadRatings: async () => {
    try {
      const stored = await SecureStore.getItem('verifier_ratings');
      if (stored) {
        set({ ratings: JSON.parse(stored) });
      }
    } catch (error) {
      console.error('Failed to load ratings:', error);
    }
  },

  getRatingForVerifier: (verifierId: string) => {
    return get().ratings.find(r => r.verifierId === verifierId) || null;
  },
}));

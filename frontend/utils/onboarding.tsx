"use client";

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";

const STORAGE_KEY = "inside-athletics.onboarding";

export type OnboardingData = {
  account: {
    name: string;
    username: string;
  };
  role: {
    role: string;
    profileImage: string | null;
  };
  preferences: {
    division: string;
    association: string;
    selectedUniversities: string[];
    primarySport: string;
    program: string;
    university: string;
  };
  verification: {
    name: string;
    email: string;
    fullName: string;
    institutionEmail: string;
    school: string;
    code: string;
  };
  plan: {
    selectedPlan: string;
  };
  topicTags: {
    selectedTags: string[];
  };
  survey: {
    responses: Record<number, string>;
  };
};

const defaultOnboardingData: OnboardingData = {
  account: {
    name: "",
    username: "",
  },
  role: {
    role: "",
    profileImage: null,
  },
  preferences: {
    division: "",
    association: "",
    selectedUniversities: [],
    primarySport: "",
    program: "",
    university: "",
  },
  verification: {
    name: "",
    email: "",
    fullName: "",
    institutionEmail: "",
    school: "",
    code: "",
  },
  plan: {
    selectedPlan: "",
  },
  topicTags: {
    selectedTags: [],
  },
  survey: {
    responses: {},
  },
};

type OnboardingContextValue = {
  data: OnboardingData;
  hydrated: boolean;
  updateSection: <K extends keyof OnboardingData>(
    section: K,
    value: Partial<OnboardingData[K]>,
  ) => void;
  reset: () => void;
};

const OnboardingContext = createContext<OnboardingContextValue | null>(null);

export function OnboardingProvider({ children }: { children: ReactNode }) {
  const [data, setData] = useState<OnboardingData>(defaultOnboardingData);
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    try {
      const storedValue = window.localStorage.getItem(STORAGE_KEY);
      if (storedValue) {
        const parsed = JSON.parse(storedValue) as Partial<OnboardingData>;
        setData({
          ...defaultOnboardingData,
          ...parsed,
          account: {
            ...defaultOnboardingData.account,
            ...parsed.account,
          },
          role: {
            ...defaultOnboardingData.role,
            ...parsed.role,
          },
          preferences: {
            ...defaultOnboardingData.preferences,
            ...parsed.preferences,
          },
          verification: {
            ...defaultOnboardingData.verification,
            ...parsed.verification,
          },
          plan: {
            ...defaultOnboardingData.plan,
            ...parsed.plan,
          },
          topicTags: {
            ...defaultOnboardingData.topicTags,
            ...parsed.topicTags,
          },
          survey: {
            ...defaultOnboardingData.survey,
            ...parsed.survey,
          },
        });
      }
    } catch {}

    setHydrated(true);
  }, []);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
  }, [data, hydrated]);

  const value = useMemo<OnboardingContextValue>(
    () => ({
      data,
      hydrated,
      updateSection(section, value) {
        setData((current) => ({
          ...current,
          [section]: {
            ...current[section],
            ...value,
          },
        }));
      },
      reset() {
        setData(defaultOnboardingData);
        window.localStorage.removeItem(STORAGE_KEY);
      },
    }),
    [data, hydrated],
  );

  return (
    <OnboardingContext.Provider value={value}>
      {children}
    </OnboardingContext.Provider>
  );
}

export function useOnboarding() {
  const context = useContext(OnboardingContext);

  if (!context) {
    throw new Error("useOnboarding must be used within an OnboardingProvider");
  }

  return context;
}

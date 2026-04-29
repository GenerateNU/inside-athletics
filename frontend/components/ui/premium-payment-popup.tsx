"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";
import { X } from "lucide-react";
import { usePostApiV1CheckoutSessions } from "@/api/hooks";
import { useSession } from "@/utils/SessionContext";
import { useCurrentUser } from "@/utils/SessionContext";

const PREMIUM_PRICE_ID = process.env.NEXT_PUBLIC_STRIPE_PREMIUM_PRICE_ID ?? "";

const planOptions = [
  {
    label: "Premium Plan",
    value: "premium" as const,
    price: "$9.99/mo",
    features: ["Access to all insider content", "Exclusive athlete insights", "Priority support"],
    gradient: "bg-[#7F8C2D]",
    card: "bg-[#E9F4A5]",
    inner: "border-[#7F8C2D] bg-[#FCFDF1]",
    divider: "border-[#7F8C2D]",
    text: "text-black",
  },
  {
    label: "Standard Plan",
    value: "free" as const,
    price: "$0/mo",
    features: ["Public posts only", "Basic college info", "Community access"],
    gradient: "bg-[#D4E94B]",
    card: "bg-[#FCFDF1]",
    inner: "border-[#D4E94B] bg-white",
    divider: "border-[#D4E94B]",
    text: "text-gray-600",
  },
] as const;

interface PremiumPaymentPopupProps {
  onClose: () => void;
}

export default function PremiumPaymentPopup({ onClose }: PremiumPaymentPopupProps) {
  const [selectedPlan, setSelectedPlan] = useState<"premium" | "free" | "">("");
  const [error, setError] = useState<string | null>(null);

  const session = useSession();
  const currentUser = useCurrentUser();

  const { mutate: createSession, isPending } = usePostApiV1CheckoutSessions({
    client: {
      headers: session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined,
    },
  });

  function handleContinue() {
    if (selectedPlan === "free") {
      onClose();
      return;
    }

    if (!currentUser?.id) {
      setError("You must be signed in to subscribe.");
      return;
    }

    if (!PREMIUM_PRICE_ID) {
      setError("Subscription is not configured. Please try again later.");
      return;
    }

    setError(null);

    createSession(
      {
        data: {
          user_id: currentUser.id,
          price_id: PREMIUM_PRICE_ID,
          quantity: 1,
          success_url: `${window.location.origin}/insidercontent?checkout=success`,
          cancel_url: `${window.location.origin}/insidercontent`,
        },
      },
      {
        onSuccess: (response) => {
          if (response.url) {
            window.location.href = response.url;
          }
        },
        onError: () => {
          setError("Failed to start checkout. Please try again.");
        },
      },
    );
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="relative flex w-full max-w-3xl max-h-[90vh] flex-col rounded-md bg-white shadow-lg mx-4">
        <div className="flex justify-end p-4 shrink-0">
          <button
            onClick={onClose}
            className="flex h-8 w-8 items-center justify-center rounded-full bg-gray-100 hover:bg-gray-200 transition-colors"
            aria-label="Close"
          >
            <X size={16} />
          </button>
        </div>

        <div className="overflow-y-auto px-8 pb-8 -mt-4 space-y-6">
          <div className="space-y-2 text-center">
            <h1 className="text-4xl font-bold text-[#001F3E]">Choose Your Plan</h1>
            <p className="text-sm text-gray-600">Pick the plan you want to start with.</p>
          </div>

          <div className="grid grid-cols-2 gap-5">
            {planOptions.map((plan) => {
              const isSelected = selectedPlan === plan.value;
              return (
                <div
                  key={plan.value}
                  className={`rounded-xl p-[2px] transition-all ${plan.gradient} ${
                    isSelected ? "-translate-y-0.5 shadow-lg" : ""
                  }`}
                >
                  <Button
                    type="button"
                    variant="outline"
                    className={`flex h-full min-h-72 w-full min-w-0 flex-col items-start rounded-[calc(0.75rem-2px)] border-transparent px-0 py-0 text-left text-black whitespace-normal ${plan.card}`}
                    onClick={() => setSelectedPlan(plan.value)}
                  >
                    <div className="w-full px-4 py-4">
                      <div className={`mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border ${plan.inner} px-3 py-3 text-sm font-semibold`}>
                        {plan.label}
                      </div>
                    </div>
                    {[plan.price, ...plan.features].map((line, i) => (
                      <div key={i} className="w-full">
                        <div className="w-full px-4">
                          <div className={`mx-3 border-t ${plan.divider}`} />
                        </div>
                        <div className={`w-full px-7 py-4 text-sm ${plan.text}`}>{line}</div>
                      </div>
                    ))}
                  </Button>
                </div>
              );
            })}
          </div>

          {error && (
            <p className="text-center text-sm text-red-500">{error}</p>
          )}

          <Button
            type="button"
            className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
            disabled={!selectedPlan || isPending}
            onClick={handleContinue}
          >
            {isPending ? (
              <span className="flex items-center gap-2">
                <Spinner className="size-4" /> Redirecting to checkout...
              </span>
            ) : (
              "Continue"
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}

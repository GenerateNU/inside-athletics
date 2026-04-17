"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { X } from "lucide-react";

type Step = "plan" | "billing";

type BillingForm = {
  cardholderName: string;
  cardNumber: string;
  expiryDate: string;
  cvc: string;
  contactEmail: string;
  contactPhone: string;
  addressLine1: string;
  addressLine2: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
};

const requiredFields: Array<keyof BillingForm> = [
  "cardholderName",
  "cardNumber",
  "expiryDate",
  "cvc",
  "contactEmail",
  "contactPhone",
  "addressLine1",
  "city",
  "state",
  "postalCode",
  "country",
];

const expiryMonths = Array.from({ length: 12 }, (_, i) =>
  String(i + 1).padStart(2, "0"),
);
const expiryYears = Array.from({ length: 12 }, (_, i) =>
  String(new Date().getFullYear() + i),
);

function parseExpiryDate(value: string) {
  const trimmed = value.trim();
  if (!trimmed) return { month: "", year: "" };
  const [rawMonth = "", rawYear = ""] = trimmed.split("/");
  const month = rawMonth.padStart(2, "0").slice(0, 2);
  if (rawYear.length === 4) return { month, year: rawYear };
  if (rawYear.length === 2) return { month, year: `20${rawYear}` };
  return { month, year: "" };
}

const planOptions = [
  { label: "Premium Plan", value: "premium", price: "$9.99/mo" },
  { label: "Standard Plan", value: "free", price: "$0/mo" },
] as const;

interface PremiumPaymentPopupProps {
  onClose: () => void;
}

export default function PremiumPaymentPopup({ onClose }: PremiumPaymentPopupProps) {
  const [step, setStep] = useState<Step>("plan");
  const [selectedPlan, setSelectedPlan] = useState("");
  const [form, setForm] = useState<BillingForm>({
    cardholderName: "",
    cardNumber: "",
    expiryDate: "",
    cvc: "",
    contactEmail: "",
    contactPhone: "",
    addressLine1: "",
    addressLine2: "",
    city: "",
    state: "",
    postalCode: "",
    country: "",
  });

  const { month: expiryMonth, year: expiryYear } = parseExpiryDate(form.expiryDate);

  const updateField = (field: keyof BillingForm, value: string) =>
    setForm((prev) => ({ ...prev, [field]: value }));

  const updateExpiryDate = ({ month = expiryMonth, year = expiryYear }: { month?: string; year?: string }) =>
    updateField("expiryDate", month && year ? `${month}/${year}` : "");

  const billingCanContinue = requiredFields.every((f) => form[f].trim().length > 0);

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
        <div className="overflow-y-auto px-8 pb-8 -mt-4">

        {step === "plan" && (
          <div className="space-y-6">
            <div className="space-y-2 text-center">
              <h1 className="text-4xl font-bold text-[#001F3E]">Choose Your Plan</h1>
              <p className="text-sm text-gray-600">
                Pick the plan you want to start with.
              </p>
            </div>

            <div className="grid grid-cols-2 gap-5">
              {planOptions.map((plan) => {
                const isSelected = selectedPlan === plan.value;
                const isPremium = plan.value === "premium";

                return (
                  <div
                    key={plan.value}
                    className={`rounded-xl p-[2px] transition-all ${
                      isPremium ? "bg-[#7F8C2D]" : "bg-[#D4E94B]"
                    } ${
                      isSelected
                        ? isPremium
                          ? "-translate-y-0.5 shadow-[4px_6px_14px_rgba(127,140,45,0.24)]"
                          : "-translate-y-0.5 shadow-[4px_6px_14px_rgba(44,100,154,0.18)]"
                        : ""
                    }`}
                  >
                    <Button
                      type="button"
                      variant="outline"
                      className={`flex h-full min-h-72 w-full min-w-0 flex-col items-start rounded-[calc(0.75rem-2px)] border-transparent px-0 py-0 text-left text-black whitespace-normal ${
                        isPremium ? "bg-[#E9F4A5]" : "bg-[#FCFDF1]"
                      }`}
                      onClick={() => setSelectedPlan(plan.value)}
                    >
                      <div className="w-full px-4 py-4">
                        {isPremium ? (
                          <div className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border-[1.5px] border-[#7F8C2D] bg-[#FCFDF1] px-3 py-3 text-sm font-semibold">
                            {plan.label}
                          </div>
                        ) : (
                          <div className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border border-[#D4E94B] bg-white px-3 py-3 text-sm font-semibold">
                            {plan.label}
                          </div>
                        )}
                      </div>
                      <div className="w-full px-4">
                        <div className={`mx-3 border-t ${isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"}`} />
                      </div>
                      <div className={`w-full px-4 py-4 text-sm ${isPremium ? "text-black" : "text-gray-600"}`}>
                        {plan.price}
                      </div>
                      <div className="w-full px-4">
                        <div className={`mx-3 border-t ${isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"}`} />
                      </div>
                      <div className={`w-full px-4 py-4 text-sm ${isPremium ? "text-black" : "text-gray-600"}`}>
                        Feature 1
                      </div>
                      <div className="w-full px-4">
                        <div className={`mx-3 border-t ${isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"}`} />
                      </div>
                      <div className={`w-full px-4 py-4 text-sm ${isPremium ? "text-black" : "text-gray-600"}`}>
                        Feature 2
                      </div>
                      <div className="w-full px-4">
                        <div className={`mx-3 border-t ${isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"}`} />
                      </div>
                      <div className={`w-full px-4 py-4 text-sm ${isPremium ? "text-black" : "text-gray-600"}`}>
                        Feature 3
                      </div>
                    </Button>
                  </div>
                );
              })}
            </div>

            <Button
              type="button"
              className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
              disabled={!selectedPlan}
              onClick={() => {
                if (selectedPlan === "free") {
                  onClose();
                } else {
                  setStep("billing");
                }
              }}
            >
              Continue
            </Button>
          </div>
        )}

        {step === "billing" && (
          <div className="space-y-6">
            <div className="space-y-2 text-center">
              <h1 className="text-4xl font-bold text-[#001F3E]">Payment Info</h1>
            </div>

            <div className="space-y-6">
              <div className="space-y-4 rounded-xl border border-gray-200 p-5">
                <h2 className="text-sm font-semibold text-[#001F3E]">Payment</h2>
                <div className="grid gap-4">
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Cardholder name</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.cardholderName}
                      onChange={(e) => updateField("cardholderName", e.target.value)}
                      placeholder="Name on card"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Card number</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.cardNumber}
                      onChange={(e) => updateField("cardNumber", e.target.value)}
                      inputMode="numeric"
                      placeholder="1234 5678 9012 3456"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Expiry date</span>
                    <div className="grid gap-4 sm:grid-cols-2">
                      <Select value={expiryMonth} onValueChange={(v) => updateExpiryDate({ month: v ?? "" })}>
                        <SelectTrigger className="h-11 w-full rounded-md border-[#3E7DBB] bg-white px-4 text-sm">
                          <SelectValue placeholder="Month" />
                        </SelectTrigger>
                        <SelectContent>
                          {expiryMonths.map((m) => (
                            <SelectItem key={m} value={m}>{m}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select value={expiryYear} onValueChange={(v) => updateExpiryDate({ year: v ?? "" })}>
                        <SelectTrigger className="h-11 w-full rounded-md border-[#3E7DBB] bg-white px-4 text-sm">
                          <SelectValue placeholder="Year" />
                        </SelectTrigger>
                        <SelectContent>
                          {expiryYears.map((y) => (
                            <SelectItem key={y} value={y}>{y}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </div>
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">CVC</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.cvc}
                      onChange={(e) => updateField("cvc", e.target.value)}
                      inputMode="numeric"
                      placeholder="123"
                    />
                  </label>
                </div>
              </div>

              <div className="space-y-4 rounded-xl border border-gray-200 p-5">
                <h2 className="text-sm font-semibold text-[#001F3E]">Billing</h2>
                <div className="grid gap-4">
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Country</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.country}
                      onChange={(e) => updateField("country", e.target.value)}
                      placeholder="Country"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Address</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.addressLine1}
                      onChange={(e) => updateField("addressLine1", e.target.value)}
                      placeholder="Street address"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">City</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.city}
                      onChange={(e) => updateField("city", e.target.value)}
                      placeholder="City"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">State</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.state}
                      onChange={(e) => updateField("state", e.target.value)}
                      placeholder="State"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">ZIP code</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.postalCode}
                      onChange={(e) => updateField("postalCode", e.target.value)}
                      placeholder="ZIP code"
                    />
                  </label>
                </div>
              </div>

              <div className="space-y-4 rounded-xl border border-gray-200 p-5">
                <h2 className="text-sm font-semibold text-[#001F3E]">Contact information</h2>
                <div className="grid gap-4">
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Contact email</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.contactEmail}
                      onChange={(e) => updateField("contactEmail", e.target.value)}
                      placeholder="name@email.com"
                      type="email"
                    />
                  </label>
                  <label className="space-y-2">
                    <span className="text-sm font-medium text-black">Contact phone</span>
                    <Input
                      className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                      value={form.contactPhone}
                      onChange={(e) => updateField("contactPhone", e.target.value)}
                      placeholder="(555) 555-5555"
                      type="tel"
                    />
                  </label>
                </div>
              </div>
            </div>

            <div className="flex gap-3">
              <Button
                type="button"
                variant="outline"
                className="h-10 flex-1 rounded-xl text-sm font-semibold"
                onClick={() => setStep("plan")}
              >
                Back
              </Button>
              <Button
                type="button"
                className="h-10 flex-1 rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
                disabled={!billingCanContinue}
                onClick={onClose}
              >
                Continue
              </Button>
            </div>
          </div>
        )}
        </div>
      </div>
    </div>
  );
}

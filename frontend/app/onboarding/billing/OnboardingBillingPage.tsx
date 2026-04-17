"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useOnboarding } from "@/utils/onboarding";

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

const expiryMonths = Array.from({ length: 12 }, (_, index) =>
  String(index + 1).padStart(2, "0"),
);

const expiryYears = Array.from({ length: 12 }, (_, index) =>
  String(new Date().getFullYear() + index),
);

function parseExpiryDate(value: string) {
  const trimmed = value.trim();

  if (!trimmed) {
    return { month: "", year: "" };
  }

  const [rawMonth = "", rawYear = ""] = trimmed.split("/");
  const month = rawMonth.padStart(2, "0").slice(0, 2);

  if (rawYear.length === 4) {
    return { month, year: rawYear };
  }

  if (rawYear.length === 2) {
    return { month, year: `20${rawYear}` };
  }

  return { month, year: "" };
}

export default function OnboardingBillingPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
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

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setForm(data.billing);
  }, [data.billing, hydrated]);

  const canContinue = requiredFields.every(
    (field) => form[field].trim().length > 0,
  );
  const { month: expiryMonth, year: expiryYear } = parseExpiryDate(
    form.expiryDate,
  );

  const updateField = (field: keyof BillingForm, value: string) => {
    setForm((current) => ({
      ...current,
      [field]: value,
    }));
  };

  const updateExpiryDate = ({
    month = expiryMonth,
    year = expiryYear,
  }: {
    month?: string;
    year?: string;
  }) => {
    updateField("expiryDate", month && year ? `${month}/${year}` : "");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-3xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">Payment Info</h1>
        </div>

        <div className="space-y-6">
          <div className="space-y-4 rounded-xl border border-gray-200 p-5">
            <h2 className="text-sm font-semibold text-[#001F3E]">Payment</h2>
            <div className="grid gap-4">
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Cardholder name
                </span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.cardholderName}
                  onChange={(event) => {
                    updateField("cardholderName", event.target.value);
                  }}
                  placeholder="Name on card"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Card number
                </span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.cardNumber}
                  onChange={(event) => {
                    updateField("cardNumber", event.target.value);
                  }}
                  inputMode="numeric"
                  placeholder="1234 5678 9012 3456"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Expiry date
                </span>
                <div className="grid gap-4 sm:grid-cols-2">
                  <Select
                    value={expiryMonth}
                    onValueChange={(value) => {
                      updateExpiryDate({ month: value ?? "" });
                    }}
                  >
                    <SelectTrigger className="h-11 w-full rounded-md border-[#3E7DBB] bg-white px-4 text-sm">
                      <SelectValue placeholder="Month" />
                    </SelectTrigger>
                    <SelectContent>
                      {expiryMonths.map((month) => (
                        <SelectItem key={month} value={month}>
                          {month}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <Select
                    value={expiryYear}
                    onValueChange={(value) => {
                      updateExpiryDate({ year: value ?? "" });
                    }}
                  >
                    <SelectTrigger className="h-11 w-full rounded-md border-[#3E7DBB] bg-white px-4 text-sm">
                      <SelectValue placeholder="Year" />
                    </SelectTrigger>
                    <SelectContent>
                      {expiryYears.map((year) => (
                        <SelectItem key={year} value={year}>
                          {year}
                        </SelectItem>
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
                  onChange={(event) => {
                    updateField("cvc", event.target.value);
                  }}
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
                  onChange={(event) => {
                    updateField("country", event.target.value);
                  }}
                  placeholder="Country"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Address
                </span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.addressLine1}
                  onChange={(event) => {
                    updateField("addressLine1", event.target.value);
                  }}
                  placeholder="Street address"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">City</span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.city}
                  onChange={(event) => {
                    updateField("city", event.target.value);
                  }}
                  placeholder="City"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">State</span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.state}
                  onChange={(event) => {
                    updateField("state", event.target.value);
                  }}
                  placeholder="State"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">ZIP code</span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.postalCode}
                  onChange={(event) => {
                    updateField("postalCode", event.target.value);
                  }}
                  placeholder="ZIP code"
                />
              </label>
            </div>
          </div>

          <div className="space-y-4 rounded-xl border border-gray-200 p-5">
            <h2 className="text-sm font-semibold text-[#001F3E]">
              Contact information
            </h2>
            <div className="grid gap-4">
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Contact email
                </span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.contactEmail}
                  onChange={(event) => {
                    updateField("contactEmail", event.target.value);
                  }}
                  placeholder="name@email.com"
                  type="email"
                />
              </label>
              <label className="space-y-2">
                <span className="text-sm font-medium text-black">
                  Contact phone
                </span>
                <Input
                  className="h-11 rounded-md border-[#3E7DBB] px-4 text-sm md:text-sm"
                  value={form.contactPhone}
                  onChange={(event) => {
                    updateField("contactPhone", event.target.value);
                  }}
                  placeholder="(555) 555-5555"
                  type="tel"
                />
              </label>
            </div>
          </div>
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={() => {
            updateSection("billing", form);
            router.push(
              role
                ? `/onboarding/topic-tags?role=${encodeURIComponent(role)}`
                : "/onboarding/topic-tags",
            );
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}

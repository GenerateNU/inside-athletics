"use client";
import { useState } from "react";
import { Plus, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Navbar } from "@/components/ui/navbar";
import { useSession } from "@/utils/SessionContext";
import {
  usePostApiV1StripeProduct,
  usePostApiV1StripePrice,
  useGetApiV1Roles,
  usePostApiV1UserByIdRoles,
  useGetApiV1UserUsernameByUsername,
  useDeleteApiV1UserByIdRoles,
} from "@/api/hooks";

type TagType =
  | "sports"
  | "divisions"
  | "athletics_performance"
  | "health_wellness"
  | "student_athlete_life"
  | "recruiting_logistics";

type MockTag = { id: string; name: string; type: TagType };

const TAG_SECTIONS: { header: string; type: TagType; max: number; group?: string }[] = [
  { header: "Sports", type: "sports", max: 1 },
  { header: "Divisions", type: "divisions", max: 3 },
  { header: "Athletics & Performance", type: "athletics_performance", max: 5, group: "Other Tags" },
  { header: "Health & Wellness", type: "health_wellness", max: 5 },
  { header: "Student Athlete Life", type: "student_athlete_life", max: 5 },
  { header: "Recruiting Logistics", type: "recruiting_logistics", max: 5 },
];

const DUMMY_TAGS: Record<TagType, MockTag[]> = {
  sports: [
    { id: "s1", name: "Basketball", type: "sports" },
    { id: "s2", name: "Football", type: "sports" },
    { id: "s3", name: "Soccer", type: "sports" },
    { id: "s4", name: "Tennis", type: "sports" },
    { id: "s5", name: "Track and Field", type: "sports" },
    { id: "s6", name: "Lacrosse", type: "sports" },
    { id: "s7", name: "Rowing", type: "sports" },
    { id: "s8", name: "Rugby", type: "sports" },
    { id: "s9", name: "Wrestling", type: "sports" },
    { id: "s10", name: "Baseball", type: "sports" },
    { id: "s11", name: "Frisbee", type: "sports" },
  ],
  divisions: [
    { id: "d1", name: "Division I", type: "divisions" },
    { id: "d2", name: "Division II", type: "divisions" },
    { id: "d3", name: "Division III", type: "divisions" },
    { id: "d4", name: "NAIA", type: "divisions" },
    { id: "d5", name: "NJCAA", type: "divisions" },
  ],
  athletics_performance: [
    { id: "ap1", name: "Coaching Style", type: "athletics_performance" },
    { id: "ap2", name: "Intensity & Competition", type: "athletics_performance" },
    { id: "ap3", name: "Player Development", type: "athletics_performance" },
    { id: "ap4", name: "Team Dynamics", type: "athletics_performance" },
  ],
  health_wellness: [
    { id: "hw1", name: "Injury & Recovery", type: "health_wellness" },
    { id: "hw2", name: "Mental Health", type: "health_wellness" },
    { id: "hw3", name: "Nutrition & Training", type: "health_wellness" },
  ],
  student_athlete_life: [
    { id: "sa1", name: "Academics & Career", type: "student_athlete_life" },
    { id: "sa2", name: "Campus & Lifestyle", type: "student_athlete_life" },
    { id: "sa3", name: "Diversity & Inclusion", type: "student_athlete_life" },
    { id: "sa4", name: "Time Management", type: "student_athlete_life" },
    { id: "sa5", name: "General", type: "student_athlete_life" },
  ],
  recruiting_logistics: [
    { id: "rl1", name: "Finances", type: "recruiting_logistics" },
    { id: "rl2", name: "NIL", type: "recruiting_logistics" },
    { id: "rl3", name: "Parents of Athletes", type: "recruiting_logistics" },
    { id: "rl4", name: "Recruiting", type: "recruiting_logistics" },
    { id: "rl5", name: "Transfer Portal", type: "recruiting_logistics" },
    { id: "rl6", name: "Walk-on Process", type: "recruiting_logistics" },
  ],
};

const DUMMY_MODERATORS = [
  { id: "u1", name: "Varun Meka", username: "@vmeka" },
  { id: "u2", name: "Mihika Chalasani", username: "@mcsquared05" },
];

const DUMMY_PRODUCTS = [
  {
    id: "prod_free",
    name: "Free",
    monthly: null,
    yearly: null,
    discountTier: "3 Free Trials",
    features: ["This day trial of Inside...", "Access to insider cont..."],
  },
  {
    id: "prod_standard",
    name: "Standard",
    monthly: 9,
    yearly: 8,
    discountTier: "Never",
    features: ["Access to insider cont...", "Interacting with other..."],
  },
];

function AddPlanModal({ onClose, onAdd }: { onClose: () => void; onAdd: () => void }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [amount, setAmount] = useState("");
  const [interval, setInterval] = useState<"day" | "week" | "month" | "year">("month");
  const [intervalCount, setIntervalCount] = useState("1");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { mutate: createProduct } = usePostApiV1StripeProduct();
  const { mutate: createPrice } = usePostApiV1StripePrice();

  const handleSubmit = () => {
    if (!name.trim() || !description.trim() || !amount.trim()) {
      setError("Please fill out all fields.");
      return;
    }
    setLoading(true);
    setError("");
    createProduct(
      { data: { name: name.trim(), description: description.trim() } },
      {
        onSuccess: (product) => {
          createPrice(
            {
              data: {
                id: product.id,
                total: Math.round(parseFloat(amount) * 100),
                interval,
                interval_count: parseInt(intervalCount),
              },
            },
            {
              onSuccess: () => { onAdd(); onClose(); },
              onError: () => { setError("Product created but failed to set price."); setLoading(false); },
            }
          );
        },
        onError: () => { setError("Failed to create plan."); setLoading(false); },
      }
    );
  };

  return (
    <div className="fixed inset-0 bg-black/30 flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl p-6 w-96 flex flex-col gap-4 shadow-lg">
        <label className="text-lg font-bold text-black">New Plan</label>
        <div className="flex flex-col gap-1">
          <label className="text-xs font-medium text-gray-600">Plan Name</label>
          <Input placeholder="e.g. Premium Plan" value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs font-medium text-gray-600">Description</label>
          <Input placeholder="e.g. Access to insider content" value={description} onChange={(e) => setDescription(e.target.value)} />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs font-medium text-gray-600">Price (USD)</label>
          <Input placeholder="e.g. 9.99" type="number" value={amount} onChange={(e) => setAmount(e.target.value)} />
        </div>
        <div className="flex gap-2">
          <div className="flex flex-col gap-1 flex-1">
            <label className="text-xs font-medium text-gray-600">Billing Interval</label>
            <select
              value={interval}
              onChange={(e) => setInterval(e.target.value as any)}
              className="border border-gray-200 rounded-lg px-3 py-2 text-sm"
            >
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
              <option value="year">Year</option>
            </select>
          </div>
          <div className="flex flex-col gap-1 w-24">
            <label className="text-xs font-medium text-gray-600">Every</label>
            <Input type="number" value={intervalCount} onChange={(e) => setIntervalCount(e.target.value)} min="1" />
          </div>
        </div>
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" onClick={onClose}>Cancel</Button>
          <Button variant="ghost" className="bg-[#2C649A] text-white hover:bg-[#245580]" onClick={handleSubmit} disabled={loading}>
            {loading ? "Creating..." : "Create"}
          </Button>
        </div>
      </div>
    </div>
  );
}

function AddTagModal({ onClose, onAdd, tagType }: { onClose: () => void; onAdd: (name: string) => void; tagType: TagType }) {
  const [name, setName] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = () => {
    if (!name.trim()) return;
    onAdd(name.trim());
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/30 flex items-center justify-center z-50">
      <div className="bg-white rounded-2xl p-6 w-80 flex flex-col gap-4 shadow-lg">
        <label className="text-lg font-bold text-black">New Tag</label>
        <Input placeholder="Tag name" value={name} onChange={(e) => setName(e.target.value)} />
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" onClick={onClose}>Cancel</Button>
          <Button variant="ghost" className="bg-[#2C649A] text-white hover:bg-[#245580]" onClick={handleSubmit}>Add</Button>
        </div>
      </div>
    </div>
  );
}

function TagSectionRow({ header, type }: { header: string; type: TagType }) {
  const [showModal, setShowModal] = useState(false);
  const [tags, setTags] = useState<MockTag[]>(DUMMY_TAGS[type]);

  const handleDelete = (id: string) => setTags((prev) => prev.filter((t) => t.id !== id));
  const handleAdd = (name: string) => setTags((prev) => [...prev, { id: `${type}-${Date.now()}`, name, type }]);

  return (
    <div className="mt-4">
      {showModal && <AddTagModal onClose={() => setShowModal(false)} onAdd={handleAdd} tagType={type} />}
      <label className="block text-sm font-bold text-black mb-1">{header}</label>
      <div className="flex flex-wrap gap-2">
        {tags.map((tag) => (
          <div key={tag.id} className="p-[0.5px]">
            <Button
              variant="ghost"
              onClick={() => handleDelete(tag.id)}
              className="rounded-lg border border-[#7F8C2D] bg-[#FCFDF1] flex items-center gap-2 px-2 py-1 font-light text-sm"
            >
              <Trash2 size={14} />
              {tag.name}
            </Button>
          </div>
        ))}
        <div className="p-[0.5px]">
          <Button
            variant="ghost"
            onClick={() => setShowModal(true)}
            className="rounded-lg border border-[#7F8C2D] bg-[#FCFDF1] flex items-center gap-2 px-2 py-1 font-light text-sm"
          >
            <Plus size={14} /> New Tag
          </Button>
        </div>
      </div>
    </div>
  );
}

function StripeProductCard({ product }: { product: typeof DUMMY_PRODUCTS[0] }) {
  return (
    <div className="border border-gray-200 rounded-xl p-4 flex flex-col gap-2 min-w-[180px]">
      <p className="font-bold text-sm text-[#001225]">{product.name}</p>
      {product.monthly !== null && <p className="text-xs text-gray-500">Monthly: ${product.monthly}</p>}
      {product.yearly !== null && <p className="text-xs text-gray-500">Yearly: ${product.yearly}</p>}
      <p className="text-xs text-gray-400">Discount: {product.discountTier}</p>
      {product.features.map((f, i) => (
        <p key={i} className="text-xs text-gray-500">Feature {i + 1}: {f}</p>
      ))}
    </div>
  );
}

function ModeratorAccessSection({ session }: { session: ReturnType<typeof useSession> }) {
  const [usernameInput, setUsernameInput] = useState("");
  const [searchUsername, setSearchUsername] = useState("");
  const [error, setError] = useState("");
  const [moderators, setModerators] = useState(DUMMY_MODERATORS);
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token ? { Authorization: `Bearer ${session.access_token}` } : undefined;

  const { data: rolesData } = useGetApiV1Roles({}, { query: { enabled }, client: { headers: authHeaders } });
  const moderatorRole = (rolesData as any)?.roles?.find((r: any) => r.name === "moderator");

  const { data: foundUser, isLoading: searching } = useGetApiV1UserUsernameByUsername(
    searchUsername,
    { query: { enabled: !!searchUsername }, client: { headers: authHeaders } }
  );

  const { mutate: assignRole } = usePostApiV1UserByIdRoles();
  const { mutate: removeRole } = useDeleteApiV1UserByIdRoles();

  const handleAdd = () => {
    if (!foundUser) { setError("User not found."); return; }
    if (!moderatorRole) { setError("Moderator role not found."); return; }
    assignRole(
      { id: (foundUser as any).id, data: { role_id: moderatorRole.id } },
      {
        onSuccess: () => {
          setModerators((prev) => [...prev, {
            id: (foundUser as any).id,
            name: `${(foundUser as any).first_name} ${(foundUser as any).last_name}`,
            username: `@${(foundUser as any).username}`,
          }]);
          setUsernameInput("");
          setSearchUsername("");
          setError("");
        },
        onError: () => setError("Failed to assign moderator role."),
      }
    );
  };

  return (
    <section className="bg-white rounded-2xl p-6 flex flex-col gap-4 border border-gray-200">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-bold text-[#001225]">Moderator Access</h2>
        <Button
          variant="ghost"
          className="bg-[#2C649A] text-white hover:bg-[#245580] rounded-full px-4 py-1 text-sm"
          onClick={handleAdd}
          disabled={!foundUser || searching}
        >
          Add New Moderator
        </Button>
      </div>
      <div className="flex gap-2 items-center">
        <Input
          placeholder="Search username"
          value={usernameInput}
          onChange={(e) => setUsernameInput(e.target.value)}
          className="max-w-xs"
        />
        <Button
          variant="ghost"
          className="border border-gray-200 rounded-lg px-3 py-2 text-sm"
          onClick={() => { setSearchUsername(usernameInput); setError(""); }}
        >
          Search
        </Button>
      </div>
      {searching && <p className="text-sm text-gray-400">Searching...</p>}
      {foundUser && !searching && (
        <div className="flex items-center gap-3 border border-gray-100 rounded-xl px-4 py-3">
          <div className="w-8 h-8 rounded-full bg-zinc-200 shrink-0" />
          <div>
            <p className="text-sm font-medium text-[#001225]">
              {(foundUser as any).first_name} {(foundUser as any).last_name}
            </p>
            <p className="text-xs text-gray-400">@{(foundUser as any).username}</p>
          </div>
        </div>
      )}
      {error && <p className="text-red-500 text-sm">{error}</p>}
      <div className="flex flex-col gap-2">
        {moderators.map((mod) => (
          <div key={mod.id} className="flex items-center justify-between border border-gray-100 rounded-xl px-4 py-3">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-zinc-200 shrink-0" />
              <div>
                <p className="text-sm font-medium text-[#001225]">{mod.name}</p>
                <p className="text-xs text-gray-400">{mod.username}</p>
              </div>
            </div>
            <Button
              variant="ghost"
              className="text-sm text-red-400 hover:text-red-600"
              onClick={() => {
                if (!moderatorRole) return;
                removeRole(
                  { id: mod.id, data: { role_id: moderatorRole.id } },
                  {
                    onSuccess: () => setModerators((prev) => prev.filter((m) => m.id !== mod.id)),
                    onError: () => setError("Failed to remove moderator."),
                  }
                );
              }}
            >
              Remove Moderator Status
            </Button>
          </div>
        ))}
      </div>
    </section>
  );
}

export default function SettingsPage() {
  const session = useSession();
  const [showAddPlanModal, setShowAddPlanModal] = useState(false);
  const otherTagsIndex = TAG_SECTIONS.findIndex((s) => s.group === "Other Tags");

  return (
    <div className="min-h-screen bg-zinc-50">
      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        {showAddPlanModal && (
          <AddPlanModal
            onClose={() => setShowAddPlanModal(false)}
            onAdd={() => console.log("plan created")}
          />
        )}
        <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
          <div className="flex w-full max-w-4xl flex-col gap-8">
            <h1 className="text-3xl font-bold text-[#001225]">Settings</h1>

            <section className="bg-white rounded-2xl p-6 flex flex-col gap-4 border border-gray-200">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-bold text-[#001225]">Subscription Options</h2>
                <Button variant="ghost" className="bg-[#2C649A] text-white hover:bg-[#245580] rounded-full px-4 py-1 text-sm">Save</Button>
              </div>
              <div className="flex flex-wrap gap-4">
                {DUMMY_PRODUCTS.map((product) => (
                  <StripeProductCard key={product.id} product={product} />
                ))}
                <div
                  onClick={() => setShowAddPlanModal(true)}
                  className="border border-dashed border-gray-300 rounded-xl p-4 flex flex-col items-center justify-center min-w-[180px] cursor-pointer hover:bg-gray-50 gap-2"
                >
                  <Plus size={20} className="text-gray-400" />
                  <p className="text-xs text-gray-400">Add New Plan</p>
                </div>
              </div>
            </section>

            <ModeratorAccessSection session={session} />

            <section className="bg-white rounded-2xl p-6 flex flex-col gap-4 border border-gray-200">
              <div className="flex justify-between items-center">
                <h2 className="text-lg font-bold text-[#001225]">Tags</h2>
                <Button variant="ghost" className="bg-[#2C649A] text-white hover:bg-[#245580] rounded-full px-4 py-1 text-sm">Save</Button>
              </div>
              {TAG_SECTIONS.map(({ header, type }, i) => (
                <div key={type}>
                  {i === otherTagsIndex && (
                    <div className="mt-4">
                      <h3 className="text-base font-bold text-[#001225]">Other Tags</h3>
                      <p className="text-xs text-[#001225]">Select Max 5</p>
                    </div>
                  )}
                  <TagSectionRow header={header} type={type} />
                  {i < otherTagsIndex && <hr className="border-t border-gray-200 mt-4" />}
                </div>
              ))}
            </section>
          </div>
        </main>
      </div>
    </div>
  );
}
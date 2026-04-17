import { ProfilePageClient } from "@/components/profile/ProfilePageClient";

type Props = {
  params: Promise<{
    id: string;
  }>;
};

export default async function OtherProfilePage({ params }: Props) {
  const { id } = await params;
  return <ProfilePageClient profileUserId={id} />;
}

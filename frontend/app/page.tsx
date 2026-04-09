import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import CreatePostPopup from "@/components/ui/create-post-popup";
import { CiUser } from "react-icons/ci";

export default function Page() {
  return (
    <div className="min-h-screen flex flex-col items-center gap-4 text-center p-6">
      <p className="text-6xl">🐸</p>
      <h1 className="text-4xl font-bold">Welcome to Inside Athletics</h1>
      <CreatePostPopup></CreatePostPopup>
      <p className="text-muted-foreground">
        Under Construction! but here are some components:
      </p>

      <div className="flex flex-row gap-5">
        <Avatar> </Avatar>
        <Avatar>
          <AvatarFallback>
            <CiUser strokeWidth={1.3} />
          </AvatarFallback>
        </Avatar>
      </div>
      <div className="flex flex-row gap-5">
        <Button variant="outline"> Click here</Button>
      </div>
    </div>
  );
}

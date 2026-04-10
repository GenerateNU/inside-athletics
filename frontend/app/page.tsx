import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Navbar } from "@/components/ui/navbar";
import { RatingPanel } from "@/components/ui/rating-panel";
import { CiUser } from "react-icons/ci";

const demoCollegeId = "014d2c09-4023-445d-9779-66aff4824245";

export default function Page() {
  return (
    <div className="min-h-screen bg-zinc-50">
      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
          <div className="flex w-full max-w-5xl flex-col items-center gap-10">
            <div className="flex flex-col items-center gap-4 text-center">
              <p className="text-6xl">🐸</p>
              <h1 className="text-4xl font-bold">Welcome to Inside Athletics</h1>
              <p className="text-muted-foreground">
                Under Construction! but here are some components:
              </p>

              <div className="flex flex-row gap-5">
                <Avatar />
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
            <RatingPanel collegeId={demoCollegeId} />
          </div>
        </main>
      </div>
    </div>
  );
}

import { useQuery } from "@tanstack/react-query";
import { Button } from "~/components/ui/button";
import { Badge } from "~/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { ArrowLeft } from "lucide-react";

export default function TmuxSessionDetailsPage() {
  const {
    data: session,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ["tmux-session", id],
    queryFn: () => fetchTmuxSessionDetails(id),
  });

  if (isLoading) {
    return (
      <main className="container mx-auto py-8 px-4">
        <div className="flex items-center justify-center p-8">
          <p className="text-muted-foreground">Loading session details...</p>
        </div>
      </main>
    );
  }

  if (isError || !session) {
    return (
      <main className="container mx-auto py-8 px-4">
        <Button
          variant="ghost"
          className="mb-6"
          onClick={() => router.push("/tmux")}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Sessions
        </Button>
        <div className="flex flex-col items-center justify-center p-8 gap-4">
          <p className="text-destructive">Failed to load session details</p>
          {error && (
            <p className="text-sm text-muted-foreground">{String(error)}</p>
          )}
        </div>
      </main>
    );
  }

  return (
    <main className="container mx-auto py-8 px-4">
      <Button
        variant="ghost"
        className="mb-6"
        onClick={() => router.push("/tmux")}
      >
        <ArrowLeft className="mr-2 h-4 w-4" />
        Back to Sessions
      </Button>

      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <h1 className="text-3xl font-bold">Tmux Session Details</h1>
          <Badge variant={session.status === "live" ? "default" : "secondary"}>
            {session.status}
          </Badge>
        </div>
        <p className="text-muted-foreground font-mono text-sm">{session.id}</p>
      </div>

      <div className="grid gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Session Information</CardTitle>
            <CardDescription>
              Basic details about this Tmux session
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Session ID
                </p>
                <p className="font-mono text-sm">{session.id}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Status
                </p>
                <p className="capitalize">{session.status}</p>
              </div>
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Number of Panes
                </p>
                <p>{session.numberOfPanes}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Commands and Output</CardTitle>
            <CardDescription>
              Recent commands executed in this session
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {session.commands.map((item, index) => (
              <div key={index} className="space-y-2">
                <div>
                  <p className="text-sm font-medium text-muted-foreground mb-1">
                    Command {index + 1}
                  </p>
                  <div className="rounded-md bg-muted p-3">
                    <code className="text-sm font-mono">{item.command}</code>
                  </div>
                </div>
                <div>
                  <p className="text-sm font-medium text-muted-foreground mb-1">
                    Output
                  </p>
                  <div className="rounded-md bg-muted p-3">
                    <pre className="text-sm font-mono whitespace-pre-wrap text-muted-foreground">
                      {item.output}
                    </pre>
                  </div>
                </div>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

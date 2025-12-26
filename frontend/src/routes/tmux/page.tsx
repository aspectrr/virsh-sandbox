import { createFileRoute } from "@tanstack/react-router";
import { TmuxSessionTable } from "~/components/tmux-session-table";

export const Route = createFileRoute()({
  component: TmuxSessionTable,
});

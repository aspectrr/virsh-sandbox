import { createFileRoute } from "@tanstack/react-router";
import { VMTable } from "~/components/vm-table";

export const Route = createFileRoute("/")({
  component: VMTable,
});

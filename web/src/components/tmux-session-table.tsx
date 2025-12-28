import * as React from "react";
import {
  useReactTable,
  getCoreRowModel,
  getSortedRowModel,
  flexRender,
  type ColumnDef,
  type SortingState,
} from "@tanstack/react-table";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { Button } from "~/components/ui/button";
import { Badge } from "~/components/ui/badge";
import { useGetApiV1TmuxSessions } from "~/tmux-client/tmux/tmux";
import { type TmuxClientInternalTypesSessionInfo } from "~/tmux-client/model";

export function TmuxSessionTable() {
  const [sorting, setSorting] = React.useState<SortingState>([]);

  const { data: sessions = [], isLoading, isError } = useGetApiV1TmuxSessions();

  const columns: ColumnDef<TmuxClientInternalTypesSessionInfo>[] = [
    {
      accessorKey: "id",
      header: "UUID",
      cell: ({ row }) => (
        <div className="font-mono text-sm">{row.getValue("id")}</div>
      ),
    },
    {
      accessorKey: "status",
      header: "Status",
      cell: ({ row }) => {
        const status = row.getValue("status") as string;
        return (
          <Badge variant={status === "live" ? "default" : "secondary"}>
            {status}
          </Badge>
        );
      },
    },
    {
      accessorKey: "vmCloneId",
      header: "VM Clone ID",
      cell: ({ row }) => (
        <div className="font-mono text-sm text-muted-foreground">
          {row.getValue("vmCloneId")}
        </div>
      ),
    },
    {
      id: "actions",
      header: "Actions",
      cell: ({ row }) => {
        return (
          <Button
            size="sm"
            onClick={() => {
              router.push(`/tmux/${row.original.id}`);
            }}
          >
            View Details
          </Button>
        );
      },
    },
  ];

  const table = useReactTable({
    data: sessions as TmuxClientInternalTypesSessionInfo[],
    columns,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    onSortingChange: setSorting,
    state: {
      sorting,
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center p-8">
        <p className="text-muted-foreground">Loading Tmux sessions...</p>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex items-center justify-center p-8">
        <p className="text-destructive">Failed to load Tmux sessions</p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border bg-card">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow
                key={row.id}
                data-state={row.getIsSelected() && "selected"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No Tmux sessions found.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}

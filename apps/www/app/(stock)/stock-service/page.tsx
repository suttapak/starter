"use client";

import { TeamCard } from "@/components/team-card";
import { NewTeamDialog } from "@/components/new-team-dialog";
import { useGetTeamMe } from "@/hooks/use-team";
import { Skeleton } from "@/components/ui/skeleton";

export default function TeamsPage() {
  const { data, isLoading, isError, error } = useGetTeamMe();
  const teams = data?.data.data;

  if (isLoading) return <Loading />;

  if (isError) throw error;

  return (
    <div className="container mx-auto py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Teams</h1>
          <p className="text-muted-foreground">
            Manage your teams and collaborate with others
          </p>
        </div>
        <NewTeamDialog />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {teams?.map((team) => (
          <TeamCard
            key={team.id}
            description={team.description}
            membersCount={0}
            name={team.name}
            username={team.username}
          />
        ))}
      </div>
      {!teams ||
        (teams.length === 0 && (
          <div className="text-center py-12">
            <p className="text-muted-foreground">
              No teams found. Create a new team to get started.
            </p>
          </div>
        ))}
    </div>
  );
}
const Loading = () => {
  return (
    <div className="container mx-auto py-8">
      <div className="flex items-center justify-between mb-8">
        <div className="flex-col flex-1 flex gap-2">
          <Skeleton className="h-8 w-28" />
          <Skeleton className="h-4 w-full" />
        </div>
        <NewTeamDialog />
      </div>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Skeleton className="h-40" />
        <Skeleton className="h-40" />
        <Skeleton className="h-40" />
      </div>
    </div>
  );
};

// ...existing code...
import { ExternalLink, GitFork, Star } from "lucide-react";
import { Card, CardContent } from "@/components/card";
import { Badge } from "@/components/badge";

interface ResultCardProps {
  repo_url: string;
  similarity: number;
  stars: number;
  forks: number;
  keywords: string[];
  description: string;
  token_similarity: number;
  metrics_similarity: number;
  ast_similarity: number;
}

const ResultCard = ({
  repo_url,
  similarity,
  stars,
  forks,
  keywords,
  description,
  token_similarity,
  metrics_similarity,
  ast_similarity,
}: ResultCardProps) => {
  const getsimilarityColor = (score: number) => {
    const s = Number.isFinite(score) ? score : 0;
    if (s >= 80) return "text-destructive";
    if (s >= 60) return "text-orange-400";
    return "text-primary";
  };

  const fmt = (v?: number) => {
    return Number.isFinite(v as number) ? `${Math.round(v as number)}%` : "—";
  };

  const safeUrl = typeof repo_url === "string" && repo_url.trim() !== "" ? repo_url.trim() : undefined;
  const displayName = safeUrl ? safeUrl.split("/").slice(-2).join("/") : "unknown/unknown";

  return (
    <Card className="group hover:shadow-lg hover:shadow-primary/5 transition-all duration-300 hover:-translate-y-1 bg-gradient-to-br from-card to-card/50 border-border/50">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <a
              href={safeUrl ?? "#"}
              target={safeUrl ? "_blank" : undefined}
              rel={safeUrl ? "noopener noreferrer" : undefined}
              className="text-lg font-semibold text-foreground hover:text-primary transition-colors flex items-center gap-2 group/link"
            >
              {displayName}
              <ExternalLink className="w-4 h-4 opacity-0 group-hover/link:opacity-100 transition-opacity" />
            </a>
            <p className="text-sm text-muted-foreground mt-1 line-clamp-2">{description}</p>
          </div>

          <div className="ml-4 text-right">
            <div className={`text-4xl font-bold ${getsimilarityColor(similarity)}`}>
              {Math.round(Number.isFinite(similarity) ? similarity : 0)}%
            </div>
            <p className="text-xs text-muted-foreground">overall similarity</p>

            <div className="mt-3 flex gap-2 justify-end">
              <Badge variant="secondary" className="text-xs px-2 py-1">
                Token: {fmt(token_similarity)}
              </Badge>
              <Badge variant="secondary" className="text-xs px-2 py-1">
                Metrics: {fmt(metrics_similarity)}
              </Badge>
              <Badge variant="secondary" className="text-xs px-2 py-1">
                AST: {fmt(ast_similarity)}
              </Badge>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-4 mb-4 text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Star className="w-4 h-4" />
            <span>{stars.toLocaleString()}</span>
          </div>
          <div className="flex items-center gap-1">
            <GitFork className="w-4 h-4" />
            <span>{forks.toLocaleString()}</span>
          </div>
        </div>

        <div className="flex flex-wrap gap-2">
          {keywords.map((Keyword, index) => (
            <Badge
              key={index}
              variant="secondary"
              className="bg-secondary/50 hover:bg-secondary text-xs"
            >
              {Keyword}
            </Badge>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};

export default ResultCard;
// ...existing code...
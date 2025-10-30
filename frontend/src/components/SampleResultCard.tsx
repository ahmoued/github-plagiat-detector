import { ExternalLink, GitFork, Star } from "lucide-react";
import { Card, CardContent } from "@/components/card";
import { Badge } from "@/components/badge";

interface ResultCardProps {
  repoUrl: string;
  similarity: number;
  stars: number;
  forks: number;
  keywords: string[];
  description: string;
}

const ResultCard = ({ repoUrl, similarity, stars, forks, keywords, description }: ResultCardProps) => {
  const getSimilarityColor = (score: number) => {
    if (score >= 80) return "text-destructive";
    if (score >= 60) return "text-orange-400";
    return "text-primary";
  };

  return (
    <Card className="group hover:shadow-lg hover:shadow-primary/5 transition-all duration-300 hover:-translate-y-1 bg-gradient-to-br from-card to-card/50 border-border/50">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <a
              href={repoUrl}
              target="_blank"
              rel="noopener noreferrer"
              className="text-lg font-semibold text-foreground hover:text-primary transition-colors flex items-center gap-2 group/link"
            >
              {repoUrl.split('/').slice(-2).join('/')}
              <ExternalLink className="w-4 h-4 opacity-0 group-hover/link:opacity-100 transition-opacity" />
            </a>
            <p className="text-sm text-muted-foreground mt-1 line-clamp-2">{description}</p>
          </div>
          
          <div className="ml-4 text-right">
            <div className={`text-3xl font-bold ${getSimilarityColor(similarity)}`}>
              {similarity}%
            </div>
            <p className="text-xs text-muted-foreground">similarity</p>
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
          {keywords.map((keyword, index) => (
            <Badge
              key={index}
              variant="secondary"
              className="bg-secondary/50 hover:bg-secondary text-xs"
            >
              {keyword}
            </Badge>
          ))}
        </div>
      </CardContent>
    </Card>
  );
};

export default ResultCard;

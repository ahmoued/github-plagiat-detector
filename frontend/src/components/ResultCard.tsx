import { Card, CardContent } from "@/components/card";

interface ResultCardProps {
  repoUrl: string;
  similarity: number;
}

const ResultCard = ({ repoUrl, similarity }: ResultCardProps) => {
  const getSimilarityColor = (score: number) => {
    if (score >= 80) return "text-destructive";
    if (score >= 60) return "text-orange-400";
    return "text-primary";
  };

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <a
            href={repoUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="text-sm hover:underline truncate flex-1"
          >
            {repoUrl.split('/').slice(-2).join('/')}
          </a>
          
          <div className={`text-2xl font-bold ml-4 ${getSimilarityColor(similarity)}`}>
            {similarity}%
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default ResultCard;

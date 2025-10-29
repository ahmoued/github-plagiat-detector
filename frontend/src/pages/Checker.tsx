import { useState } from "react";
import { Button } from "@/components/button";
import { Input } from "@/components/input";
import { toast } from "sonner";
import ResultCard from "@/components/ResultCard";
import { Link } from "react-router-dom";

const Checker = () => {
  const [repoUrl, setRepoUrl] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [results, setResults] = useState<any[]>([]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!repoUrl.trim()) {
      toast.error("Please enter a URL");
      return;
    }

    if (!repoUrl.includes("github.com")) {
      toast.error("Please enter a valid GitHub URL");
      return;
    }

    setIsLoading(true);

    setTimeout(() => {
      const mockResults = [
        {
          repoUrl: "https://github.com/example/similar-project",
          similarity: 87,
        },
        {
          repoUrl: "https://github.com/demo/code-analyzer",
          similarity: 76,
        },
        {
          repoUrl: "https://github.com/test/repo-checker",
          similarity: 65,
        },
      ];

      setResults(mockResults);
      setIsLoading(false);
      toast.success("Done!");
    }, 2000);
  };

  return (
    <div className="min-h-screen flex flex-col p-4">
      <div className="mb-8 text-center">
        <Link
          to="/"
          className="text-sm text-muted-foreground hover:text-foreground"
        >
          ← Back
        </Link>
      </div>

      <main className="flex-1 max-w-2xl mx-auto w-full space-y-8">
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-bold">Check Repository</h1>
          <p className="text-muted-foreground">Enter GitHub repository URL</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            type="text"
            placeholder="https://github.com/username/repository"
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            className="h-12"
            disabled={isLoading}
          />
          <Button type="submit" disabled={isLoading} className="w-full h-12">
            {isLoading ? "Checking..." : "Check"}
          </Button>
        </form>

        {isLoading && (
          <div className="text-center py-8">
            <p className="text-muted-foreground">Analyzing...</p>
          </div>
        )}

        {!isLoading && results.length > 0 && (
          <div className="space-y-4">
            <h2 className="text-xl font-bold">Results</h2>

            {results.map((result, index) => (
              <ResultCard key={index} {...result} />
            ))}
          </div>
        )}
      </main>
    </div>
  );
};

export default Checker;

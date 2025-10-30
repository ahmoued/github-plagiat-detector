import { useEffect } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import { Card, CardContent } from "@/components/card";
import { Button } from "../components/button";
import { ArrowLeft, Copy, ExternalLink, Star, GitFork } from "lucide-react";
import { toast } from "sonner";

const pct = (v: any) => {
  const n = Number(v ?? 0);
  if (!Number.isFinite(n)) return 0;
  return Math.max(0, Math.min(100, Math.round(n)));
};

const ResultDetail = () => {
  const { state } = useLocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!state || !state.result) {

        const raw = sessionStorage.getItem("plagiarism:lastResults");
      const rawRepo = sessionStorage.getItem("plagiarism:lastRepo");
      if (!raw || !rawRepo) {
        navigate("/checker");
        return;
      }

      
    }
  }, [state, navigate]);

  if (!state || !state.result) {

    return (
      <div className="min-h-screen flex flex-col bg-gradient-to-b from-background to-background/95">
        <Navbar />
        <main className="flex-1 container mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-16">
          <div className="max-w-4xl mx-auto">
            <Card>
              <CardContent className="p-8 text-center">
                <h2 className="text-xl font-semibold mb-2">
                  No result selected
                </h2>
                <p className="text-sm text-muted-foreground">
                  Please go back to the checker and select a result. Your last
                  run's results are preserved, so you can return to them.
                </p>
                <div className="mt-6">
                  <Button onClick={() => navigate("/checker")}>
                    Back to results
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </main>
        <Footer />
      </div>
    );
  }

  const { result, overall } = state as any;

  const token = pct(result.token_similarity ?? result.tokenSimilarity ?? 0);
  const metrics = pct(
    result.metrics_similarity ?? result.metricsSimilarity ?? 0
  );
  const ast = pct(result.ast_similarity ?? result.astSimilarity ?? 0);
  const stars = result.stars ?? 0;
  const forks = result.forks ?? 0;
  const description = result.description ?? "";
  const repo_url = result.repo_url ?? result.repoUrl ?? "";

  const copyRepo = async () => {
    try {
      await navigator.clipboard.writeText(repo_url);
      toast.success("Repository URL copied");
    } catch (e) {
      toast.error("Could not copy URL");
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-b from-background to-background/95">
      <Navbar />
      <main className="flex-1 container mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-16">
        <div className="max-w-4xl mx-auto">
          <Card>
            <CardContent className="p-8">
              <div className="flex items-center justify-between gap-4">
                <div className="flex items-center gap-3">
                  <Button variant="ghost" onClick={() => navigate(-1)}>
                    <ArrowLeft className="w-4 h-4" />
                  </Button>
                  <div>
                    <h2 className="text-xl font-semibold">
                      {(repo_url || "").split("/").slice(-2).join("/")}
                    </h2>
                    <p className="text-xs text-muted-foreground">
                      Overall similarity
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <div className="text-right">
                    <div className="text-5xl font-extrabold mb-1">
                      {overall}%
                    </div>
                    <div className="text-sm text-muted-foreground">
                      ⭐ {stars} • Forks: {forks}
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      onClick={copyRepo}
                      title="Copy repo URL"
                    >
                      <Copy className="w-4 h-4" />
                    </Button>
                    <a
                      href={repo_url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      <Button variant="ghost" title="Open in GitHub">
                        <ExternalLink className="w-4 h-4" />
                      </Button>
                    </a>
                  </div>
                </div>
              </div>

              <div className="mt-6 space-y-6">
                <div>
                  <div className="flex justify-between mb-1">
                    <div className="text-sm font-medium text-gray-700">
                      Token similarity
                    </div>
                    <div className="text-sm font-medium text-gray-700">
                      {token}%
                    </div>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className="bg-primary h-3 rounded-full transition-all duration-700"
                      style={{ width: `${token}%` }}
                    />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between mb-1">
                    <div className="text-sm font-medium text-gray-700">
                      Metrics similarity
                    </div>
                    <div className="text-sm font-medium text-gray-700">
                      {metrics}%
                    </div>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className="bg-secondary h-3 rounded-full transition-all duration-700"
                      style={{ width: `${metrics}%` }}
                    />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between mb-1">
                    <div className="text-sm font-medium text-gray-700">
                      AST similarity
                    </div>
                    <div className="text-sm font-medium text-gray-700">
                      {ast}%
                    </div>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-3">
                    <div
                      className="bg-destructive h-3 rounded-full transition-all duration-700"
                      style={{ width: `${ast}%` }}
                    />
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-semibold mb-1">Description</h3>
                  <p className="text-sm text-muted-foreground">{description}</p>
                </div>

                <div className="flex gap-2 mt-4">
                  {(result.keywords ?? []).map((k: string, i: number) => (
                    <span
                      key={i}
                      className="px-2 py-1 bg-muted rounded text-xs"
                    >
                      {k}
                    </span>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  );
};

export default ResultDetail;

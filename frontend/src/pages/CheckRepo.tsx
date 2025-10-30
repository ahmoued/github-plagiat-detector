import { useState, useEffect } from "react";
import { Search, Loader2 } from "lucide-react";
import { Button } from "../components/button";
import { Input } from "../components/input";
import { toast } from "sonner";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer";
import ResultCard from "../components/SampleResultCard";
import axios from "axios";
import { useNavigate } from "react-router-dom";

type Result = {
  repo: string;
  repo_url: string;
  token_similarity: number;
  metrics_similarity: number;
  ast_similarity: number;
  stars: number;
  forks: number;
  keywords: string[];
  description: string;
};

const Checker = () => {
  const [repo_url, setrepo_url] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [results, setResults] = useState<Result[]>([]);
  const navigate = useNavigate();


  useEffect(() => {
    try {
      const raw = sessionStorage.getItem("plagiarism:lastResults");
      const rawRepo = sessionStorage.getItem("plagiarism:lastRepo");
      if (raw) {
        const parsed = JSON.parse(raw) as Result[];
        setResults(parsed);
      }
      if (rawRepo) setrepo_url(rawRepo);
    } catch (e) {

      console.warn("failed to restore session results", e);
    }
  }, []);

  const api = axios.create({
    baseURL: "http://localhost:8080",
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!repo_url.trim()) {
      toast.error("Please enter a GitHub repository URL");
      return;
    }

    if (!repo_url.includes("github.com")) {
      toast.error("Please enter a valid GitHub repository URL");
      return;
    }
    if (repo_url.includes(".git")) {
      toast.error("Please remove the .git extension");
      return;
    }

    setIsLoading(true);
    try {
      const res = await api.post("/compare", { repo_url: repo_url });
      if (res.status !== 200) {
        console.log("error checking the repo", res.statusText);
      }
      console.log("front");
      console.log(res.data.results);
      console.log("debugging");
      setResults(res.data.results);
      try {
        sessionStorage.setItem(
          "plagiarism:lastResults",
          JSON.stringify(res.data.results)
        );
        sessionStorage.setItem("plagiarism:lastRepo", repo_url);
      } catch (e) {
        console.warn("failed to persist results", e);
      }
    } catch (err) {
      console.log(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-b from-background to-background/95">
      <Navbar />

      <main className="flex-1 container mx-auto px-4 sm:px-6 lg:px-8 pt-24 pb-16">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12 animate-fade-in">
            <h1 className="text-4xl sm:text-5xl font-bold mb-4 bg-gradient-to-r from-primary via-accent to-primary bg-clip-text text-transparent">
              repository Plagiarism Checker
            </h1>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Enter a GitHub repository URL to analyze for potential code
              similarities and duplicates
            </p>
          </div>

          <form onSubmit={handleSubmit} className="mb-12 animate-fade-in-up">
            <div className="flex flex-col sm:flex-row gap-4 p-6 rounded-2xl bg-card/50 backdrop-blur-sm border border-border shadow-lg">
              <Input
                type="text"
                placeholder="https://github.com/username/repository"
                value={repo_url}
                onChange={(e) => setrepo_url(e.target.value)}
                className="flex-1 h-12 bg-background/50 border-border focus:border-primary transition-colors"
                disabled={isLoading}
              />
              <Button
                type="submit"
                disabled={isLoading}
                className="h-12 px-8 bg-primary hover:bg-primary/90 text-primary-foreground font-semibold shadow-lg shadow-primary/20 hover:shadow-xl hover:shadow-primary/30 transition-all"
              >
                {isLoading ? (
                  <>
                    <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                    Analyzing...
                  </>
                ) : (
                  <>
                    <Search className="w-5 h-5 mr-2" />
                    Check repository
                  </>
                )}
              </Button>
            </div>
          </form>

          {isLoading && (
            <div className="text-center py-12 animate-fade-in">
              <div className="inline-block p-4 rounded-full bg-primary/10 animate-glow mb-4">
                <Loader2 className="w-12 h-12 text-primary animate-spin" />
              </div>
              <p className="text-muted-foreground">
                Scanning repositories and analyzing code patterns...
              </p>
            </div>
          )}

          {!isLoading && results.length > 0 && (
            <div className="space-y-6 animate-fade-in">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-2xl font-bold text-foreground">
                  Analysis Results
                </h2>
                <p className="text-sm text-muted-foreground">
                  Found {results.length} potential matches
                </p>
              </div>

              <div className="grid gap-6">
                {results.map((result, index) => {
                  const tok = Number(result.token_similarity ?? 0);
                  const met = Number(result.metrics_similarity ?? 0);
                  const ast = Number(result.ast_similarity ?? -1);

                  let overall = 0;
                  if (Number.isFinite(ast) && ast > 0) {
                    overall = (ast + tok + met) / 3;
                  } else {
                    overall = (tok + met) / 2;
                  }
                  if (!Number.isFinite(overall)) overall = 0;

                  return (
                    <div
                      key={index}
                      className="animate-fade-in-up"
                      style={{ animationDelay: `${index * 0.1}s` }}
                    >
                      <ResultCard
                        similarity={Math.round(overall)}
                        repo_url={result.repo_url}
                        metrics_similarity={Math.round(
                          result.metrics_similarity ?? 0
                        )}
                        ast_similarity={Math.round(result.ast_similarity ?? 0)}
                        token_similarity={Math.round(
                          result.token_similarity ?? 0
                        )}
                        stars={result.stars ?? 0}
                        forks={result.forks ?? 0}
                        description={result.description ?? ""}
                        keywords={result.keywords ?? []}
                        onClick={() =>
                          navigate("/result", {
                            state: {
                              result,
                              overall: Math.round(overall),
                            },
                          })
                        }
                      />
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {!isLoading && results.length === 0 && repo_url && (
            <div className="text-center py-12 animate-fade-in">
              <div className="inline-block p-4 rounded-full bg-secondary/50 mb-4">
                <Search className="w-12 h-12 text-muted-foreground" />
              </div>
              <p className="text-muted-foreground">
                No results yet. Enter a repository URL and click "Check
                repository"
              </p>
            </div>
          )}
        </div>
      </main>

      <Footer />
    </div>
  );
};

export default Checker;

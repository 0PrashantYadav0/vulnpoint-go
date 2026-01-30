import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Textarea } from "@/components/ui/textarea";
import { Sparkles, Loader2, Bot, Wand2 } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";

interface AIAssistantSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onGenerate: (prompt: string) => Promise<void>;
  isGenerating: boolean;
}

const AIAssistantSheet = ({
  open,
  onOpenChange,
  onGenerate,
  isGenerating,
}: AIAssistantSheetProps) => {
  const [prompt, setPrompt] = useState("");

  const handleGenerate = async () => {
    if (!prompt.trim()) return;
    await onGenerate(prompt);
  };

  const handleQuickPrompt = (text: string) => {
    setPrompt(text);
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-md flex flex-col h-full">
        <SheetHeader className="space-y-4 pb-4 border-b">
          <SheetTitle className="flex items-center gap-2 text-xl">
            <Bot className="w-6 h-6 text-indigo-500" />
            AI Architect
          </SheetTitle>
          <SheetDescription className="text-base text-gray-500 dark:text-gray-400">
            Describe your security workflow needs, and I'll architect the perfect pipeline for you.
          </SheetDescription>
        </SheetHeader>

        <div className="flex-1 py-6 space-y-6">
          <div className="space-y-3">
            <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              Your Requirements
            </label>
            <Textarea
              placeholder="e.g. I need to scan a GitHub repo for secrets every 6 hours, run a dependency check, and if anything critical is found, create a high-priority Jira ticket and slack the security team..."
              value={prompt}
              onChange={(e) => setPrompt(e.target.value)}
              className="min-h-[150px] resize-none text-base p-4 focus-visible:ring-indigo-500 bg-slate-50 dark:bg-slate-900 border-slate-200 dark:border-slate-800"
            />
            <p className="text-xs text-muted-foreground text-right">
              Be as specific as you like about tools and conditions.
            </p>
          </div>

          <div className="space-y-3">
             <div className="flex items-center gap-2">
                <Wand2 className="w-3.5 h-3.5 text-indigo-500" />
                <label className="text-sm font-medium">Quick Start Templates</label>
             </div>
             <ScrollArea className="h-[200px] rounded-md border p-4 bg-slate-50 dark:bg-slate-900/50">
                <div className="space-y-2">
                  <button 
                    onClick={() => handleQuickPrompt("Scan my github repo for secrets and vulnerable configurations using Semgrep, then create a GitHub issue for findings.")}
                    className="w-full text-left text-sm p-3 rounded-md bg-white dark:bg-slate-800 border hover:border-indigo-500 hover:shadow-sm transition-all duration-200"
                  >
                    <div className="font-medium text-slate-900 dark:text-slate-100 mb-1">Secrets & SAST Pipeline</div>
                    <div className="text-slate-500 dark:text-slate-400 text-xs line-clamp-2">Full code scan with Gitleaks + Semgrep &gt; GitHub Issue</div>
                  </button>

                  <button 
                     onClick={() => handleQuickPrompt("Run a full OWASP scan on my website every day at midnight. If vulnerabilities are found, email security@company.com.")}
                     className="w-full text-left text-sm p-3 rounded-md bg-white dark:bg-slate-800 border hover:border-indigo-500 hover:shadow-sm transition-all duration-200"
                   >
                    <div className="font-medium text-slate-900 dark:text-slate-100 mb-1">Daily Web Audit</div>
                    <div className="text-slate-500 dark:text-slate-400 text-xs line-clamp-2">Scheduled Nikto/OWASP scan &gt; Email Report</div>
                  </button>
                  
                   <button 
                     onClick={() => handleQuickPrompt("Check for new CVEs in my dependencies using Trivy. If critical, use auto-fix to create a PR updating the package.")}
                     className="w-full text-left text-sm p-3 rounded-md bg-white dark:bg-slate-800 border hover:border-indigo-500 hover:shadow-sm transition-all duration-200"
                   >
                    <div className="font-medium text-slate-900 dark:text-slate-100 mb-1">Auto-Patch Dependencies</div>
                    <div className="text-slate-500 dark:text-slate-400 text-xs line-clamp-2">Trivy Scan &gt; Auto-Fix PR for critical CVEs</div>
                  </button>
                </div>
             </ScrollArea>
          </div>
        </div>

        <SheetFooter className="border-t pt-4 sm:justify-between items-center bg-slate-50 -mx-6 -mb-6 px-6 py-4 dark:bg-slate-900">
           <Button variant="ghost" onClick={() => onOpenChange(false)}>
             Cancel
           </Button>
          <Button
            onClick={handleGenerate}
            disabled={!prompt.trim() || isGenerating}
            className="w-full sm:w-auto gap-2 bg-indigo-600 hover:bg-indigo-700 shadow-lg shadow-indigo-500/20 transition-all duration-200"
          >
            {isGenerating ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin" />
                Architecting Workflow...
              </>
            ) : (
              <>
                <Sparkles className="w-4 h-4" />
                Generate Workflow
              </>
            )}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};

export default AIAssistantSheet;

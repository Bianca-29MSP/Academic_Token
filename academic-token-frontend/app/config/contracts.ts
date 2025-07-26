// Contract addresses for Academic Token
export const CONTRACTS = {
  equivalence: 'academic14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sva26tx',
  progress: 'academic1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqqg0cfq',
  degree: 'academic17p9rzwnnfxcjp32un9ug7yhhzgtkhvl9jfksztgw5uh69wac2pgscn3aqy',
  schedule: 'academic1ghd753shjuwexxywmgs4xz7x2q732vcnkm6h2pyv9s6ah3hylvrqy9uyq9',
  academicNFT: 'academic1eyfccmjm6732k7wp4p6gdjwhxjwsvje44j0hfx8nkgrm8fs7vqfsljrz85'
};

// Contract query messages
export const EQUIVALENCE_QUERIES = {
  getState: { get_state: {} },
  analyzeEquivalence: (sourceSubjectId: string, targetSubjectId: string) => ({
    analyze_equivalence: {
      source_subject_id: sourceSubjectId,
      target_subject_id: targetSubjectId
    }
  }),
  getAnalysisResult: (analysisId: string) => ({
    get_analysis_result: {
      analysis_id: analysisId
    }
  })
};

// Contract execute messages
export const EQUIVALENCE_EXECUTES = {
  requestAnalysis: (sourceSubjectId: string, targetSubjectId: string, studentId: string) => ({
    request_analysis: {
      source_subject_id: sourceSubjectId,
      target_subject_id: targetSubjectId,
      student_id: studentId
    }
  })
};

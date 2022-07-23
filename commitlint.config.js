// Configuration of commitlint to check commit message guidelines
module.exports = {
    extends: ['@commitlint/config-conventional'],
    parserPreset: 'conventional-changelog-conventionalcommits',
    rules: {
        'subject-max-length': [2, 'always', 50],
        'scope-enum': [2, 'always', [
            'deps', // Changes done on anything dependency related
            'devops', // Changes done on technical processes
            'api', // Changes in /api/ directory
            'comp', // Changes in /component/ directory
            'int', // Changes in /internal/ directory
            'core' // Changes on files in project root
        ]]
    }
};

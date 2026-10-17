import { mdiCheckAll, mdiContentCopy, mdiContentSaveAll } from '@mdi/js';
import Icon from '@mdi/react';
import { saveAs } from 'file-saver';
import { useCallback, useState } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';

import { setClipboardText } from '../../helpers/utils/setClipboardText';

interface ICodeBlock {
    readonly language?: string;
    readonly fileName?: string;
    readonly children: string;
    readonly showLineNumbers?: boolean;
}

export const CodeBlock = ({ fileName, language, children, showLineNumbers }: ICodeBlock) => {
    const [isCopied, setIsCopied] = useState(false);

    const handleCopy = useCallback(async (text: string) => {
        await setClipboardText(text);
        setIsCopied(true);
        await new Promise((resolve) => setTimeout(resolve, 2000));
        setIsCopied(false);
    }, []);

    const handleDownload = useCallback(
        (text: string) => {
            const blob = new Blob([text], {
                type: 'text/plain;charset=utf-8'
            });
            saveAs(blob, fileName);
        },
        [fileName]
    );

    return (
        <div className="app-code-block">
            <div className="app-code-toolbar">
                <span className="app-code-filename">{fileName ?? 'Code'}</span>
                <div className="app-code-actions">
                    <div
                        className="app-icon-button"
                        onClick={() => {
                            handleCopy(children);
                        }}
                    >
                        <Icon
                            className="app-code-icon"
                            path={isCopied ? mdiCheckAll : mdiContentCopy}
                            size={0.8}
                        />
                    </div>
                    {fileName?.length && (
                        <div
                            className="app-icon-button"
                            onClick={() => {
                                handleDownload(children);
                            }}
                        >
                            <Icon className="app-code-icon" path={mdiContentSaveAll} size={0.8} />
                        </div>
                    )}
                </div>
            </div>
            <SyntaxHighlighter
                language={language}
                useInlineStyles={false}
                codeTagProps={{ className: language ? `language-${language}` : undefined }}
                showLineNumbers={showLineNumbers}
            >
                {children}
            </SyntaxHighlighter>
        </div>
    );
};

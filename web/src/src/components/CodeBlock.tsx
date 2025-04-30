import { mdiCheckAll, mdiContentCopy, mdiContentSaveAll } from '@mdi/js';
import Icon from '@mdi/react';
import { saveAs } from 'file-saver';
import { useCallback, useState } from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import theme from 'react-syntax-highlighter/dist/esm/styles/prism/atom-dark';

import { setClipboardText } from '../helpers/utils/setClipboardText';

interface ICodeBlock {
    readonly language?: string;
    readonly fileName?: string;
    readonly children: string;
}

export const CodeBlock = ({ fileName, language, children }: ICodeBlock) => {
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
        <div className="rounded-lg bg-gray-700 p-2">
            <div className="flex items-center justify-between px-3">
                <span className="font-mono text-sm text-gray-300">{fileName ?? 'Code'}</span>
                <div className="flex space-x-3">
                    <div
                        className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                        onClick={() => {
                            handleCopy(children);
                        }}
                    >
                        <Icon
                            className="flex-shrink-0 text-gray-300"
                            path={isCopied ? mdiCheckAll : mdiContentCopy}
                            size={0.8}
                        />
                    </div>
                    {fileName?.length && (
                        <div
                            className="cursor-pointer opacity-60 transition-all hover:opacity-100"
                            onClick={() => {
                                handleDownload(children);
                            }}
                        >
                            <Icon
                                className="flex-shrink-0 text-gray-300"
                                path={mdiContentSaveAll}
                                size={0.8}
                            />
                        </div>
                    )}
                </div>
            </div>
            <SyntaxHighlighter language={language} style={theme}>
                {children}
            </SyntaxHighlighter>
        </div>
    );
};
